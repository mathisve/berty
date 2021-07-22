package cmd

import (
	"context"
	"errors"
	"fmt"
	"infratesting/aws"
	"infratesting/logging"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"infratesting/config"
	"infratesting/daemon/grpc/daemon"
	"infratesting/testing"
	"log"
	"math/rand"
)

var (
	testCmd = &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx := context.Background()
			t := time.Now().Unix()

			c, err := loadConfig()
			if err != nil {
				return logging.LogErr(err)
			}

			availablePeers,  err := testing.GetAllEligiblePeers(aws.Ec2TagType, []string{config.NodeTypePeer})
			if err != nil {
				return logging.LogErr(err)
			}

			availableRepl, err := testing.GetAllEligiblePeers(aws.Ec2TagType, []string{config.NodeTypeReplication})
			if err != nil {
				return logging.LogErr(err)
			}

			log.Printf("PEERS: %d desired, %d available\n", c.CountPeers(), len(availablePeers))
			if len(availablePeers) != c.CountPeers() {
				return logging.LogErr(errors.New("amount of available peers doesn't match desired peers"))
			}

			log.Printf("REPLICATION: %d desired, %d available\n", c.CountRepl(), len(availableRepl))
			if len(availableRepl) != c.CountRepl() {
				return logging.LogErr(errors.New("amount of available replication peers doesn't match desired replication peers"))
			}

			// temporary group object
			type group struct {
				Name string
				Pk []byte
				Peers []*testing.Peer
				Tests []config.Test
			}

			var groups map[string]group
			groups = make(map[string]group, c.GetAmountOfGroups())

			// assign peers and tests to group
			for i := range availablePeers {
				availablePeers[i].MatchNodeToPeer(c)
				for g := range availablePeers[i].ConfigGroups {
					group := groups[availablePeers[i].ConfigGroups[g].Name]

					group.Name = availablePeers[i].ConfigGroups[g].Name
					group.Peers = append(group.Peers, &availablePeers[i])

					if len(group.Tests) == 0 {
						group.Tests = availablePeers[i].ConfigGroups[g].Tests
					}

					group.Tests = append(group.Tests)

					groups[availablePeers[i].ConfigGroups[g].Name] = group
				}
			}

			// match replication server to node
			for i := range availableRepl {
				availableRepl[i].MatchNodeToPeer(c)
				for g := range availableRepl[i].ConfigGroups {
					p := groups[availableRepl[i].ConfigGroups[g].Name].Peers[rand.Intn(len(groups[availableRepl[i].ConfigGroups[g].Name].Peers))]
					_, err = p.P.AddReplication(ctx, &daemon.AddReplication_Request{
						GroupName: availableRepl[g].ConfigGroups[g].Name,
						TokenIp: availableRepl[i].Ip,
					})

					if err != nil {
						return logging.LogErr(err)
					}
				}
			}

			// create and join groups

			for i := range availablePeers {
				for g := range availablePeers[i].ConfigGroups {
					if len(groups[availablePeers[i].ConfigGroups[g].Name].Pk) == 0 {
						// group doesn't exist yet
						invite, err := availablePeers[i].P.CreateInvite(ctx, &daemon.CreateInvite_Request{GroupName: availablePeers[i].ConfigGroups[g].Name})
						if err != nil {
							return logging.LogErr(err)
						}

						gr := groups[availablePeers[i].ConfigGroups[g].Name]
						gr.Pk = invite.Invite
						groups[availablePeers[i].ConfigGroups[g].Name] = gr

						_, err = availablePeers[i].P.StartReceiveMessage(ctx, &daemon.StartReceiveMessage_Request{
							GroupName: availablePeers[i].ConfigGroups[g].Name,
						})
						if err != nil {
							return logging.LogErr(err)
						}

					} else {
						// group already exists, join it
						_, err := availablePeers[i].P.JoinGroup(ctx, &daemon.JoinGroup_Request{
							GroupName: availablePeers[i].ConfigGroups[g].Name,
							Invite:    groups[availablePeers[i].ConfigGroups[g].Name].Pk,
						})
						if err != nil {
							return logging.LogErr(err)
						}

						_, err = availablePeers[i].P.StartReceiveMessage(ctx, &daemon.StartReceiveMessage_Request{
							GroupName: availablePeers[i].ConfigGroups[g].Name,
						})
						if err != nil {
							return logging.LogErr(err)
						}
					}
				}
			}

			var groupArray []group

			for key := range groups {
				groupArray = append(groupArray, groups[key])
			}

			var startTestWG sync.WaitGroup
			var finishTestWG sync.WaitGroup


			// iterate over groups
			for g := range groupArray {
				// iterate over tests in groups
				for j := range groupArray[g].Tests {
					// iterate over peers in group
					for k := range groupArray[g].Peers {
						groupIndex := g
						testIndex := j
						peerIndex := k
						go func(startTestWG, finishTestWG *sync.WaitGroup) {
							startTestWG.Add(1)
							finishTestWG.Add(1)
							// create new test with correct variables
							_, err = groupArray[groupIndex].Peers[peerIndex].P.NewTest(ctx, &daemon.NewTest_Request{
								GroupName: groupArray[groupIndex].Name,
								TestN:  int64(j),
								Type:      groupArray[groupIndex].Tests[testIndex].TypeInternal,
								Size:      int64(groupArray[groupIndex].Tests[testIndex].SizeInternal),
								Interval:  int64(groupArray[groupIndex].Tests[testIndex].IntervalInternal),
								Amount:    int64(groupArray[groupIndex].Tests[testIndex].AmountInternal),
							})
							logging.Log(fmt.Sprintf("added test: %v", j))

							if err != nil {
								log.Println(err)
							}

							startTestWG.Done()

							fmt.Println("waiting")
							startTestWG.Wait()

							// start said test
							_, err = groupArray[groupIndex].Peers[peerIndex].P.StartTest(ctx, &daemon.StartTest_Request{
								GroupName: groupArray[groupIndex].Name,
								TestN:  int64(j),
							})
							if err != nil {
								log.Println(err)
							}

							finishTestWG.Done()
						}(&startTestWG, &finishTestWG)
					}
					finishTestWG.Wait()
					time.Sleep(time.Second * 4)
					log.Printf("Started test %v in group %v on\n", j, groupArray[g].Name)
				}

				var wg sync.WaitGroup

				for j := range groupArray[g].Peers {
					for k := range groupArray[g].Tests {

						wg.Add(1)

						testIndex := k
						peerIndex := j
						go func(wg *sync.WaitGroup) {
							for {
								isRunning, err := groupArray[g].Peers[peerIndex].P.IsTestRunning(ctx, &daemon.IsTestRunning_Request{
									GroupName: groupArray[g].Name,
									TestN: int64(testIndex),
								})
								if err != nil {
									logging.Log(err)
								} else {
									if isRunning.TestIsRunning == false {
										break
									}
									time.Sleep(time.Second * 3)
								}
							}
							wg.Done()
						}( &wg)
					}
				}

				fmt.Println("waiting for tests to finish ...")
				wg.Wait()
				fmt.Println("all tests are finished!")

 			}



 			allNodes, err := testing.GetAllEligiblePeers(aws.Ec2TagType, []string{
 				config.NodeTypePeer,
 				config.NodeTypeReplication,
 				config.NodeTypeRDVP,
				config.NodeTypeRelay,
				config.NodeTypeBootstrap,
 			})

			for k := range allNodes{
				resp, err := allNodes[k].P.UploadLogs(ctx, &daemon.UploadLogs_Request{
					Folder: strconv.FormatInt(t, 10),
					Name:   strings.ReplaceAll(allNodes[k].Tags[aws.Ec2TagName], ".", "-"),
				})
				if err != nil {
					return logging.LogErr(err)
				}

				logging.Log(fmt.Sprintf("uploaded : %v\n", resp.UploadCount))
			}

			return nil
		},
	}
	)
