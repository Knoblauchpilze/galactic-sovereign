<script lang="ts">
	import { FlexContainer, StyledButton, StyledTitle } from '@totocorpsoftwareinc/frontend-toolkit';

	let { data } = $props();

	const labelStyle = 'font-bold px-2 text-right';
	const fieldStyle = 'text-[#2e5a88]';
</script>

<FlexContainer>
	<StyledTitle text="User details" />

	<table class="text-secondary">
		<tbody>
			<tr>
				<td class={labelStyle}>ID:</td>
				<td class={fieldStyle}>{data.user.id}</td>
			</tr>
			<tr>
				<td class={labelStyle}>E-mail:</td>
				<td class={fieldStyle}>{data.user.email}</td>
			</tr>
			<tr>
				<td class={labelStyle}>Password:</td>
				<td class={fieldStyle}>{data.user.password}</td>
			</tr>
			<tr>
				<td class={labelStyle}>API key:</td>
				<td class={fieldStyle}>{data.apiKey}</td>
			</tr>
			<tr>
				<td class={labelStyle}>Member since:</td>
				<td class={fieldStyle}>{data.user.createdAt}</td>
			</tr>
		</tbody>
	</table>

	<form method="POST" action="?/logout">
		<StyledButton text="Logout" />
	</form>
</FlexContainer>
