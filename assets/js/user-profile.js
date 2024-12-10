class UserProfile extends HTMLElement {
  async connectedCallback() {
    const splitCookie = document.cookie.split(";");
    const crsfTokenFull = splitCookie.find((cookie) =>
      cookie.includes("csrf_token")
    );
    const crsfToken = crsfTokenFull.split("=").at(1);

    if (!crsfToken) {
      return;
    }

    const userData = await fetch("/auth/me", {
      headers: {
        "X-CSRF-TOKEN": crsfToken,
      },
    });

    const json = await userData.json();

    this.innerHTML = `
          <img src="${json.profile_picture}" alt="${json.username}'s avatar" class="profile-picture">
        `;
  }
}

customElements.define("user-profile", UserProfile);
