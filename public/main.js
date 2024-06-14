export function getUser() {
	const cookies = new Map(document.cookie.split(";").map(record => record.split('=')))
	const jwt = cookies.get('jwt')
	if (!jwt) {
		return null
	}
	const [, payloadB64] = jwt.split('.')
	const payload = JSON.parse(atob(payloadB64))
	return payload
}

export class User {
	/**
	* @param {string} uuid 
	* @param {string} name 
	* @param {string} email 
	* @param {string} role 
	* @param {string | undefined} profilePicture 
	*/
	constructor(uuid, username, email, role, profilePicture) {
		this.uuid = uuid
		this.name = username
		this.email = email
		this.role = role
		this.profilePicture = profilePicture
	}

	getProfilePicture() {
		return this.profilePicture || 'https://placehold.it/300/300'
	}
}

export async function whoami() {
	try {
		const r = await fetch('/api/v0/whoami')
		if (r.ok) {
			const user = await r.json()
			return new User(user.uuid, user.username, user.email, user.role, user.profilePicture)
		} else {
			throw new Error(r.statusText)
		}
	} catch (e) {
		throw new Error('Unable to fetch the current user')
	}
}
