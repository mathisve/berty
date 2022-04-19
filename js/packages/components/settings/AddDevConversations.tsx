import React from 'react'
import { View, ScrollView } from 'react-native'
import { Layout } from '@ui-kitten/components'
import { useTranslation } from 'react-i18next'

import { useStyles } from '@berty/styles'
import { globals } from '@berty/config'
import { useThemeColor } from '@berty/store/hooks'

import { ButtonSetting } from '../shared-components/SettingsButtons'
import { ScreenFC, useNavigation } from '@berty/navigation'

type ValueOf<T> = T[keyof T]

type ValueType = (
	| ValueOf<typeof globals.berty.contacts>
	| ValueOf<typeof globals.berty.conversations>
) & { kind: string }

const uncap = (s: string) => s[0].toLowerCase() + s.slice(1)

const Button: React.FC<ValueType> = ({ kind: rawKind, name, link }) => {
	const { t } = useTranslation()
	const colors = useThemeColor()
	const navigation = useNavigation()

	const kind = uncap(rawKind)

	const text = t('settings.add-dev-conversations.tag-' + kind)
	console.log('kind', kind)
	let state
	switch (kind) {
		case 'bot':
			state = {
				value: text,
				color: colors['background-header'],
				bgColor: colors['positive-asset'],
			}
			break
		case 'contact':
			state = {
				value: text,
				color: colors['background-header'],
				bgColor: colors['positive-asset'],
			}
			break
		case 'conversation':
			state = {
				value: text,
				color: colors['background-header'],
				bgColor: colors['positive-asset'],
			}
			break
		default:
			state = {
				value: 'unknown',
				color: colors['secondary-text'],
				bgColor: colors['main-background'],
			}
	}

	return (
		<ButtonSetting
			name={t('settings.add-dev-conversations.add', { name: name })}
			icon='book-outline'
			iconSize={30}
			iconColor={colors['alt-secondary-background-header']}
			actionIcon={null}
			state={state}
			onPress={() => {
				navigation.navigate('Modals.ManageDeepLink', { type: 'link', value: link })
			}}
		/>
	)
}

const BodyAddContactList = () => {
	const [{ padding, flex, margin }] = useStyles()

	return (
		<View style={[padding.medium, flex.tiny, margin.bottom.small]}>
			{Object.values(globals.berty.contacts).map(value => {
				return <Button key={value.link} {...value} />
			})}
			{Object.values(globals.berty.conversations).map(value => {
				return <Button key={value.link} {...value} kind='Conversation' />
			})}
		</View>
	)
}

export const AddDevConversations: ScreenFC<'Settings.AddDevConversations'> = () => {
	const colors = useThemeColor()

	return (
		<Layout style={{ flex: 1, backgroundColor: colors['main-background'] }}>
			<ScrollView bounces={false}>
				<BodyAddContactList />
			</ScrollView>
		</Layout>
	)
}
