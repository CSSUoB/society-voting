export const API = {
    AUTH_LOGIN: `auth/login`,
    AUTH_LOGOUT: `auth/logout`,

    ME: `/api/me`,
    ME_NAME: `/api/me/name`,

    ELECTION: `api/election`,

    ADMIN_ELECTION: `api/admin/election`,
    ADMIN_USER: `api/admin/user`,
    ADMIN_USER_DELETE: `/api/admin/user/delete`,
} as const;