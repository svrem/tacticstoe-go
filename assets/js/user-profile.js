class UserProfile extends HTMLElement {
  async connectedCallback() {
    const split_cookie = document.cookie.split(";");
    const crsf_token_full = split_cookie.find((cookie) =>
      cookie.includes("csrf_token")
    );
    const crsf_token = crsf_token_full.split("=").at(1);

    const user_data = await fetch("/auth/me", {
      headers: {
        "X-CSRF-TOKEN": crsf_token,
      },
    });

    const json = await user_data.json();

    this.innerHTML = `
          <img src="${json.profile_picture}" alt="${json.username}'s avatar" class="profile-picture">
        `;
  }
}

customElements.define("user-profile", UserProfile);
