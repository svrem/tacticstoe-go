class Auth {
  constructor() {
    this.user = null;

    this.onAuthChangeCallbacks = [];

    this.loadUserData().then(() => {
      this.onAuthChangeCallbacks.forEach((cb) => cb(this.user));
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
    cb(this.user);
    this.onAuthChangeCallbacks.push(cb);
  }

  async loadUserData() {
    const userRes = await fetch("/auth/me");

    if (!userRes.ok) {
      return;
    }

    const userData = await userRes.json();

    this.user = userData;
  }
}

window.Auth = new Auth();
