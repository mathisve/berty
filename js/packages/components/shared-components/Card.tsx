import React, { Fragment } from 'react'
import {
	View,
	ViewProps,
	TouchableHighlight,
	TouchableHighlightProps,
	StyleSheet,
} from 'react-native'

import { useThemeColor } from '@berty-tech/store/hooks'

const style = StyleSheet.create({
	default: {
		borderRadius: 20,
		padding: 16,
		margin: 16,
	},
})

export const TouchableCard: React.FunctionComponent<TouchableHighlightProps> = (props) => {
	const colors = useThemeColor()

	return (
		<TouchableHighlight
			activeOpacity={0.9}
			underlayColor={colors['main-text']}
			{...props}
			style={[style.default, props.style]}
			onPress={props.onPress ? props.onPress : (): null => null}
		>
			<Fragment>{props.children}</Fragment>
		</TouchableHighlight>
	)
}

export const Card: React.FunctionComponent<ViewProps> = (props) => (
	<View {...props} style={[style.default, props.style]}>
		<Fragment>{props.children}</Fragment>
	</View>
)

export default Card
