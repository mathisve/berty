import React from 'react'
import { TouchableOpacity, View } from 'react-native'
import { CommonActions } from '@react-navigation/native'

import { useStyles } from '@berty/styles'
import { useConversation } from '@berty/hooks'
import { dispatch } from '@berty/navigation'

import { useStylesNotification, NotificationTmpLogo } from './common'
import { UnifiedText } from '../shared-components/UnifiedText'

const GroupInvitation: React.FC<any> = ({ onClose, title, message, ...props }) => {
	const [{ text }] = useStyles()
	const _styles = useStylesNotification()

	const { payload } = props?.additionalProps?.payload || {}
	const conv = useConversation(payload.contact?.conversationPublicKey)

	const handlePressConvMessage = () => {
		if (conv?.publicKey) {
			dispatch(
				CommonActions.reset({
					routes: [
						{ name: 'Main.Home' },
						{
							name: 'Chat.OneToOne',
							params: {
								convId: conv.publicKey,
							},
						},
					],
				}),
			)
		} else {
			console.warn('Notif: ContactGroupInvitation: Conversation not found or no public key')
		}
		if (typeof onClose === 'function') {
			onClose()
		}
	}

	return (
		<TouchableOpacity
			style={_styles.touchable}
			activeOpacity={conv ? 0.3 : 1}
			//underlayColor='transparent'
			onPress={handlePressConvMessage}
		>
			<View style={_styles.innerTouchable}>
				<NotificationTmpLogo />
				<View style={_styles.titleAndTextWrapper}>
					<UnifiedText numberOfLines={1} style={[text.bold]}>
						{title}
					</UnifiedText>
					<UnifiedText numberOfLines={1} ellipsizeMode='tail'>
						{message}
					</UnifiedText>
				</View>
			</View>
		</TouchableOpacity>
	)
}

export default GroupInvitation
