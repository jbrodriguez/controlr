// messages - mutations
export const IS_BUSY = 'mutation/IS_BUSY'
export const LOGOUT = 'mutation/LOGOUT'
export const LOGIN_SUCCESS = 'mutation/LOGIN_SUCCESS'
export const LOGIN_FAILURE = 'mutation/LOGIN_FAILURE'
export const REMOVE_FEEDBACK = 'mutation/REMOVE_FEEDBACK'
export const SET_SOCKET = 'mutation/SET_SOCKET'
export const SELECT_USER = 'mutation/SELECT_USER'
export const UPDATE_ALL = 'mutation/UPDATE_ALL'
export const SET_ERROR = 'mutation/SET_ERROR'

// messages - actions
export const GET_VERSION = 'action/GET_VERSION'
export const GOT_VERSION = 'action/GOT_VERSION'
export const LOGIN = 'action/LOGIN'
export const INIT = 'action/INIT'
export const SET_PERM = 'action/SET_PERM'

// socket - actions (incoming)
export const PERMS_UPDATED = 'model/USER_UPDATED'
export const REFRESHED = 'model/REFRESHED'
export const ACCESS_ERROR = 'model/ACCESS_ERROR'

// socket - actions (outgoing)
export const REFRESH = 'model/REFRESH'
export const UPDATE_PERMS = 'model/UPDATE_USER'

// routes
export const HOME = '/'
export const PERMS = '/permissions'
