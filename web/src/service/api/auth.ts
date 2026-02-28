import { request } from '../request';

export function fetchLogin(email: string, password: string) {
  return request<Api.Auth.LoginToken>({
    url: '/user/login',
    method: 'post',
    data: {
      email,
      password
    }
  });
}

/** Get user info */
export function fetchGetUserInfo() {
  return request<Api.Auth.UserInfo>({ url: '/user/info' });
}

/**
 * Refresh token
 *
 * @param refreshToken Refresh token
 */
export function fetchRefreshToken(refreshToken: string) {
  return request<Api.Auth.LoginToken>({
    url: '/auth/refreshToken',
    method: 'post',
    data: {
      refreshToken
    }
  });
}

/**
 * return custom backend error
 *
 * @param code error code
 * @param msg error message
 */
export function fetchCustomBackendError(code: string, msg: string) {
  return request({ url: '/auth/error', params: { code, msg } });
}

export function UserRegister(username: string, password: string, email: string) {
  return request<Api.Auth.registerResponse>({
    url: '/user/register',
    method: 'post',
    data: {
      email,
      username,
      password
    }
  });
}

/** Update user profile info (username, email) */
export function updateUserProfile(data: { username: string; email: string }) {
  return request({
    url: '/user/updateInfo',
    method: 'post',
    data
  });
}

/** Update user password */
export function updateUserPassword(data: { first_password: string; second_password: string }) {
  return request({
    url: '/user/resetPassword',
    method: 'post',
    data
  });
}
