class Dialog extends HTMLElement {
  constructor() {
    super();

    this.attachShadow({ mode: "open" });

    this.shadowRoot.innerHTML = `
      <style>
        .dialog {
            display: grid;
            place-items: center;
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;



            & > * {
                opacity: 0;
            }
            
            &.hidden > * {
                pointer-events: none;
                opacity: 0;
            
            }

            &.show > * {
                pointer-events: all;
                opacity: 1;
            }
        }

        .dialog__content {
            background-color: white;
            padding: 10px;
            z-index: 1000;

            margin: 1rem;

            position: relative;

            width: 80%;
            max-width: 600px;

            padding: 1rem;
            border-radius: .5rem;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.2);

            transition: opacity 0.1s;
        }

        .dialog__overlay {
            position: fixed;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background-color: rgba(0, 0, 0, 0.3);

            transition: opacity 0.1s;
        }

        .close {
            position: absolute;
            top: 0;
            right: 0;
            background: none;
            border: none;
            padding: 1rem;
            cursor: pointer;
        }
      </style>
      <div class="dialog">
        <div class="dialog__overlay"></div>
        <div class="dialog__content">
            <button id="close" class="close">
                <svg  xmlns="http://www.w3.org/2000/svg"  width="24"  height="24"  viewBox="0 0 24 24"  fill="none"  stroke="black"  stroke-width="2"  stroke-linecap="round"  stroke-linejoin="round"  class="icon icon-tabler icons-tabler-outline icon-tabler-x"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M18 6l-12 12" /><path d="M6 6l12 12" /></svg>
            </button>

            <slot></slot>
        </div>
      </div>
    `;
  }

  hide() {
    this.ariaHidden = true;
    this.tabIndex = -1;

    this.shadowRoot.querySelector(".dialog").classList.add("hidden");
    this.shadowRoot.querySelector(".dialog").classList.remove("show");

    setTimeout(() => {
      this.remove();
    }, 100);
  }

  show() {
    this.ariaHidden = false;
    this.tabIndex = 1;

    this.shadowRoot.querySelector(".dialog").classList.add("show");
    this.shadowRoot.querySelector(".dialog").classList.remove("hidden");
  }

  connectedCallback() {
    this.shadowRoot.querySelector("#close").addEventListener("click", () => {
      this.hide();
    });

    this.shadowRoot
      .querySelector(".dialog__overlay")
      .addEventListener("click", () => {
        this.hide();
      });
  }
}

customElements.define("app-dialog", Dialog);
