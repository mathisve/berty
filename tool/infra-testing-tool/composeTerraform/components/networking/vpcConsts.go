package networking

const (
	// VpcHCLTemplate is the Vpc template
	// refer to https://golang.org/pkg/text/template/ for templating syntax
	VpcHCLTemplate = `
resource "aws_vpc" "{{.Name }}" {
  cidr_block = "{{.CidrBlock }}"
  enable_dns_hostnames = true
  enable_dns_support = true
}
`
	// VpcCidrBlockDefault is the default value for CidrBlock
	VpcCidrBlockDefault = "10.0.1.0/24"

	// VpcNamePrefix is the prefix for the Vpc type
	VpcNamePrefix = "vpc"

	VpcType = "vpc"

	// VpcErrNoValidCidrBlock means the attached CidrBlock is not valid
	VpcErrNoValidCidrBlock = "vpc has no valid attribute 'CidrBlock'"
)
