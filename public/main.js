export function getUser() {
	const cookies = new Map(document.cookie.split(";").map(record => record.split('=')))
	const jwt = cookies.get('jwt')
	if (!jwt) {
		return null
	}
	const [, payloadB64] = jwt.split('.')
	return JSON.parse(atob(payloadB64))
}
