import React from 'react'
import { ScrollView, View } from 'react-native'

import { useStyles } from '@berty/styles'
import { ScreenFC } from '@berty/navigation'
import { useThemeColor } from '@berty/store'

import { ButtonSettingV2, Section } from '../shared-components'

export const Security: ScreenFC<'Settings.Security'> = () => {
	const [{}, { scaleSize }] = useStyles()
	const colors = useThemeColor()

	return (
		<View style={{ backgroundColor: colors['secondary-background'], flex: 1 }}>
			<ScrollView
				bounces={false}
				contentContainerStyle={{ paddingBottom: 12 * scaleSize }}
				showsVerticalScrollIndicator={false}
			>
				<Section>
					{/* TODO i18n */}
					<ButtonSettingV2 text='TouchID/FaceID before open account' disabled />
					<ButtonSettingV2 text='Password' disabled last />
				</Section>
			</ScrollView>
		</View>
	)
}
