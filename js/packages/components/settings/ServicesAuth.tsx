import React, { useState } from 'react'
import { View, ScrollView } from 'react-native'
import { Layout, Input } from '@ui-kitten/components'
import { useTranslation } from 'react-i18next'

import { useStyles } from '@berty-tech/styles'
import { ScreenProps, useNavigation } from '@berty-tech/navigation'
import { useMsgrContext, useThemeColor } from '@berty-tech/store/hooks'
import {
	servicesAuthViaURL,
	servicesAuthViaDefault,
	useAccountServices,
	serviceNames,
} from '@berty-tech/store/services'

import { HeaderSettings, ButtonSetting, FactionButtonSetting } from '../shared-components'

const BodyServicesAuth = () => {
	const [{ flex, padding, margin }] = useStyles()
	const colors = useThemeColor()
	const { t } = useTranslation()

	const [url, setURL] = useState('')
	const ctx = useMsgrContext()
	const accountServices = useAccountServices()

	return (
		<View style={[flex.tiny, padding.medium, margin.bottom.medium]}>
			<ButtonSetting
				name={t('settings.services-auth.operated-services-button')}
				icon='plus-circle-outline'
				iconSize={30}
				iconColor={colors['background-header']}
				alone={true}
				onPress={async () => {
					await servicesAuthViaDefault(ctx)
				}}
			/>
			<FactionButtonSetting
				name={t('settings.services-auth.register-service-button.title')}
				icon='plus-circle-outline'
				iconSize={30}
				iconColor={colors['background-header']}
				style={[margin.top.medium]}
			>
				<Input
					textContentType={'URL'}
					autoCorrect={false}
					autoCapitalize={'none'}
					value={url}
					placeholder={t('settings.services-auth.register-service-button.input-placeholder')}
					onChange={({ nativeEvent }) => {
						setURL(nativeEvent.text)
					}}
				/>
				<ButtonSetting
					name={t('settings.services-auth.register-service-button.action')}
					iconSize={30}
					iconColor={colors['background-header']}
					alone={false}
					onPress={async () => {
						await servicesAuthViaURL(ctx, url)
					}}
				/>
			</FactionButtonSetting>
			<FactionButtonSetting
				name={t('settings.services-auth.registered-services-button.title')}
				icon='cube-outline'
				iconSize={30}
				iconColor={colors['background-header']}
				style={[margin.top.medium]}
			>
				{accountServices.length === 0 ? (
					<ButtonSetting
						name={t('settings.services-auth.registered-services-button.sample-no-services')}
						disabled
						alone={false}
					/>
				) : (
					accountServices.map((t) => {
						return (
							<ButtonSetting
								key={`${t.tokenId}-${t.serviceType}`}
								name={`${
									(typeof t.serviceType === 'string' && serviceNames[t.serviceType]) ||
									t('settings.services-auth.registered-services-button.sample-unknown-service')
								}\n${t.authenticationUrl}`}
								disabled
								alone={false}
							/>
						)
					})
				)}
			</FactionButtonSetting>
		</View>
	)
}

export const ServicesAuth: React.FC<ScreenProps.Settings.ServicesAuth> = () => {
	const { goBack } = useNavigation()
	const { t } = useTranslation()
	const colors = useThemeColor()

	return (
		<Layout style={{ flex: 1, backgroundColor: colors['main-background'] }}>
			<ScrollView bounces={false}>
				<HeaderSettings title={t('settings.services-auth.title')} undo={goBack} />
				<BodyServicesAuth />
			</ScrollView>
		</Layout>
	)
}
