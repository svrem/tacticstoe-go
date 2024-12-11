class Auth {
  constructor() {
    this.user = null;

    self.onAuthChangeCallbacks = [];

    this.loadUserData().then(() => {
      self.onAuthChangeCallbacks.forEach((cb) => cb(self.user));
    });
  }

  login() {
    window.location.href = "/auth/login/google";
  }

  logout(cb) {
    cb();
  }

  isAuthenticated() {
    return this.authenticated;
  }

  onAuthChange(cb) {
    cb(self.user);
    self.onAuthChangeCallbacks.push(cb);
  }

  async loadUserData() {
    const splitCookie = document.cookie.split(";");
    const crsfTokenFull = splitCookie.find((cookie) =>
      cookie.includes("csrf_token")
    );

    if (!crsfTokenFull) {
      return;
    }

    const crsfToken = crsfTokenFull.split("=").at(1);

    if (!crsfToken) {
      return;
    }

    const userRes = await fetch("/auth/me", {
      headers: {
        "X-CSRF-TOKEN": crsfToken,
      },
    });

    if (!userRes.ok) {
      return;
    }

    const userData = await userRes.json();

    self.user = userData;
  }
}

window.Auth = new Auth();
