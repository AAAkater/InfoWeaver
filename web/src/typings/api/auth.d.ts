declare namespace Api {
  /**
   * namespace Auth
   *
   * backend api module: "auth"
   */
  namespace Auth {
    interface registerResponse {
      code: number;
      msg: string;
      data: string;
    }
    interface LoginToken {
      token: string;
      refreshToken: string;
    }

    interface UserInfo {
      id: string;
      username: string;
      email: string;
      roles: string[];
      buttons: string[];
    }
  }
}
