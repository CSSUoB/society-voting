export const API = {
	AUTH_LOGIN: `/auth/login`,
	AUTH_LOGOUT: `/auth/logout`,

	ME: `/api/me`,
	ME_NAME: `/api/me/name`,

	ELECTION: `/api/election`,
	ELECTION_STAND: `/api/election/stand`,
	ELECTION_CURRENT: `/api/election/current`,
	ELECTION_CURRENT_VOTE: `/api/election/current/vote`,
	ELECTION_SSE: `/api/election/sse`,
	ELECTION_RESULTS: `/api/election/results`,

	ADMIN_ELECTION: `/api/admin/election`,
	ADMIN_ELECTION_SSE: `/api/admin/election/sse`,
	ADMIN_ELECTION_START: `/api/admin/election/start`,
	ADMIN_ELECTION_STOP: `/api/admin/election/stop`,
	ADMIN_ELECTION_PUBLISH: `/api/admin/election/publish`,
	ADMIN_USER: `/api/admin/user`,
	ADMIN_USER_DELETE: `/api/admin/user/delete`,
	ADMIN_USER_RESTRICT: `/api/admin/user/restrict`,
} as const;
