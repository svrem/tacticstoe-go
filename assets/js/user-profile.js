window.onProfilePictureError = (event) => {
  event.target.src = "/assets/svg/user-question.svg";
};

class UserProfile extends HTMLElement {
  async connectedCallback() {
    window.Auth.onAuthChange((user) => {
      if (user) {
        this.innerHTML = `
          <img src="${user.profile_picture}" alt="${user.username}'s avatar" class="profile-picture" onerror="window.onProfilePictureError(event)" />
        `;
      } else {
        this.innerHTML = "";
      }
    });
  }
}

customElements.define("user-profile", UserProfile);
