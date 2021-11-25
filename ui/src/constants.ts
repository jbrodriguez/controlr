// messages - mutations
export const IS_BUSY = 'mutation/IS_BUSY'
export const LOGOUT = 'mutation/LOGOUT'
export const LOGIN_SUCCESS = 'mutation/LOGIN_SUCCESS'
export const LOGIN_FAILURE = 'mutation/LOGIN_FAILURE'
export const REMOVE_FEEDBACK = 'mutation/REMOVE_FEEDBACK'
export const SET_SOCKET = 'mutation/SET_SOCKET'
export const UPDATE_ALL = 'mutation/UPDATE_ALL'
export const SET_ERROR = 'mutation/SET_ERROR'
export const SET_PERM = 'mutation/SET_PERM'
export const CLEAR_ALL = 'mutation/CLEAR_ALL'
export const UPDATE_PERMS = 'mutation/UPDATE_PERMS'

// messages - actions
export const GET_VERSION = 'action/GET_VERSION'
export const GOT_VERSION = 'action/GOT_VERSION'
export const LOGIN = 'action/LOGIN'
export const INIT = 'action/INIT'
export const SAVE_CHANGES = 'action/SAVE_CHANGES'
export const GET_BARCODE = 'action/GET_BARCODE'
export const GOT_BARCODE = 'action/GOT_BARCODE'

// socket - actions (incoming)
export const USER_UPDATED = 'model/USER_UPDATED'
export const REFRESHED = 'model/REFRESHED'
export const ACCESS_ERROR = 'model/ACCESS_ERROR'

// socket - actions (outgoing)
export const REFRESH = 'model/REFRESH'
export const UPDATE_USER = 'model/UPDATE_USER'

// routes
export const HOME = '/'
export const PERMS = '/permissions'
