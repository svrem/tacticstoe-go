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
    const userRes = await fetch("/auth/me");

    if (!userRes.ok) {
      return;
    }

    const userData = await userRes.json();

    self.user = userData;
  }
}

window.Auth = new Auth();
