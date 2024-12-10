class UserProfile extends HTMLElement {
  async connectedCallback() {
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
    window.userData = userData;

    this.innerHTML = `
          <img src="${userData.profile_picture}" alt="${userData.username}'s avatar" class="profile-picture">
        `;
  }
}

customElements.define("user-profile", UserProfile);
