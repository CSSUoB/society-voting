export const API = {
	AUTH_LOGIN: `/auth/login`,
	AUTH_LOGOUT: `/auth/logout`,

	ME: `/api/me`,
	ME_NAME: `/api/me/name`,

	POLL: `/api/poll`,
	POLL_CURRENT: `/api/poll/current`,
	POLL_RESULTS: `/api/poll/results`,
	POLL_SSE: `/api/poll/sse`,
	ELECTION_STAND: `/api/election/stand`,
	ELECTION_VOTE: `/api/election/vote`,
	REFERENDUM_VOTE: `/api/referendum/vote`,

	ADMIN_POLL: `/api/admin/poll`,
	ADMIN_POLL_PUBLISH: `/api/admin/poll/publish`,
	ADMIN_POLL_SSE: `/api/admin/poll/sse`,
	ADMIN_ELECTION: `/api/admin/election`,
	ADMIN_ELECTION_START: `/api/admin/election/start`,
	ADMIN_ELECTION_STOP: `/api/admin/election/stop`,
	ADMIN_REFERENDUM: `/api/admin/referendum`,
	ADMIN_REFERENDUM_START: `/api/admin/referendum/start`,
	ADMIN_REFERENDUM_STOP: `/api/admin/referendum/stop`,
	ADMIN_USER: `/api/admin/user`,
	ADMIN_USER_DELETE: `/api/admin/user/delete`,
	ADMIN_USER_RESTRICT: `/api/admin/user/restrict`,
} as const;

export enum PollTypeId {
	ELECTION = 1,
	REFERENDUM
}

type EndpointType = "vote" | "create" | "start" | "stop";

const endpoints = {
  1: {
	vote: API.ELECTION_VOTE,
	create: API.ADMIN_ELECTION,
	start: API.ADMIN_ELECTION_START,
	stop: API.ADMIN_ELECTION_STOP,
  },
  2: {
	vote: API.REFERENDUM_VOTE,
	create: API.ADMIN_REFERENDUM,
	start: API.ADMIN_REFERENDUM_START,
	stop: API.ADMIN_REFERENDUM_STOP, 
  },
} as {[key: number]: {[key: string]: string}};

export const getEndpointForPollType = (endpointType: EndpointType, pollType: PollTypeId): string | undefined => {
	return endpoints[pollType]?.[endpointType];
};