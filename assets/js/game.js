class BoardModal extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {
    this.classList.add("game-end-message-container");
    this.classList.add("visible");

    this.innerHTML = `

        <div class="top">
          <h2>${this.getAttribute("title")}</h2>

          <button class="close" onclick="this.parentElement.parentElement.hide()">
            <svg  xmlns="http://www.w3.org/2000/svg"  width="24"  height="24"  viewBox="0 0 24 24"  fill="none"  stroke="black"  stroke-width="2"  stroke-linecap="round"  stroke-linejoin="round"  class="icon icon-tabler icons-tabler-outline icon-tabler-x"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M18 6l-12 12" /><path d="M6 6l12 12" /></svg>
          </button>
        </div>

        <p>${this.getAttribute("message")}</p> 

        <div class="bottom">
          <button onclick="openSocket()">Play Again</button>
        </div>
  `;
    console.log(this.innerHTML);
  }

  hide() {
    this.classList.remove("visible");
    this.classList.add("hidden");

    setTimeout(() => {
      this.remove();
    }, 300);
  }
}

customElements.define("board-modal", BoardModal);

class GameContainer extends HTMLElement {
  constructor() {
    super();

    this.cellCallbacks = [];

    this.classList.add("game-container");
    this.id = "game-container";

    this.innerHTML = `
      <div class="game-info">
      </div>

      <div class="game-board" id="game-board">
      </div>

      <div class="game-info">
      </div>
    `;

    this.game_board = this.querySelector("#game-board");

    this.reset();
  }

  onCellClick(callback) {
    this.cellCallbacks.push(callback);
  }

  handleCellClick = (e) => {
    const cell_index = e.target.getAttribute("data-cell");

    const x = cell_index % 4;
    const y = Math.floor(cell_index / 4);

    this.cellCallbacks.forEach((callback) => {
      callback({ x, y });
    });
  };

  showModal(title, message) {
    const boardModal = document.createElement("board-modal");
    boardModal.setAttribute("title", title);
    boardModal.setAttribute("message", message);
    this.appendChild(boardModal);
  }

  hideModal() {
    const boardModal = this.querySelector("board-modal");
    if (boardModal) boardModal.hide();
  }

  reset() {
    for (const el of this.querySelectorAll(".cell")) {
      el.remove();
    }

    this.removeAttribute("data-draw");
    this.removeAttribute("data-player-winner");
    this.removeAttribute("data-player-turn");

    this.hideModal();

    this.cellCallbacks = [];

    for (let i = 0; i < 16; i++) {
      const cell = document.createElement("button");

      cell.setAttribute("data-cell", i);
      cell.classList.add("cell");
      cell.addEventListener("click", this.handleCellClick);

      this.game_board.appendChild(cell);
    }
  }

  setGameEnd(winner, coords) {
    switch (winner) {
      case "draw":
        this.showModal("Draw!", "No one won this game.");
        this.setAttribute("data-draw", true);
        break;
      case "player":
        this.showModal("You Won!", "Your ELO has increased by 10 points!");
        this.setAttribute("data-player-winner", true);
        break;
      case "opponent":
        this.showModal("You Lost!", "Your ELO has decreased by 10 points!");
        this.setAttribute("data-player-winner", false);
        break;
    }

    let coords_index = 0;

    const i = setInterval(() => {
      const [x, y] = coords[coords_index];
      const cell_index = y * 4 + x;
      const cell = this.game_board.children[cell_index];

      cell.setAttribute("data-winner", winner === "player");

      coords_index++;

      if (coords_index >= coords.length) {
        clearInterval(i);
      }
    }, 300);
  }

  handleGameUpdate({ x, y }, state, is_active_player) {
    this.setAttribute("data-player-turn", is_active_player);

    const cell_index = y * 4 + x;
    const cell = this.game_board.children[cell_index];

    cell.setAttribute("data-state", state);

    if (state === 1) {
      cell.innerHTML =
        '<svg  xmlns="http://www.w3.org/2000/svg"  width="72"  height="72"  viewBox="0 0 24 24"  fill="none"  stroke="black"  stroke-width="2"  stroke-linecap="round"  stroke-linejoin="round"  class="icon icon-tabler icons-tabler-outline icon-tabler-x"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M18 6l-12 12" /><path d="M6 6l12 12" /></svg>';
    } else if (state === 2) {
      cell.innerHTML =
        '<svg  xmlns="http://www.w3.org/2000/svg"  width="72"  height="72"  viewBox="0 0 24 24"  fill="none"  stroke="black"  stroke-width="2"  stroke-linecap="round"  stroke-linejoin="round"  class="icon icon-tabler icons-tabler-outline icon-tabler-circle"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M12 12m-9 0a9 9 0 1 0 18 0a9 9 0 1 0 -18 0" /></svg>';
    }
  }
}

customElements.define("game-container", GameContainer);

let game_container = document.getElementById("game-container");

function openSocket() {
  game_container = document.getElementById("game-container");

  game_container.reset();

  const socket = new WebSocket("/ws");

  const game_state = {
    player_id: null,
    active_player: null,
  };

  socket.onopen = function (e) {
    console.log("Connected to server");
  };
  socket.onmessage = function (event) {
    handleWebSocketMessage(event, game_state);
  };

  game_container.onCellClick(({ x, y }) => {
    if (game_state.active_player !== game_state.player_id) {
      return;
    }

    socket.send(
      JSON.stringify({
        type: "action",
        data: {
          x,
          y,
        },
      })
    );
  });
}

function handleWebSocketMessage(event, game_state) {
  const server_message = JSON.parse(event.data);

  switch (server_message.type) {
    case "game_start":
      initializeGame(server_message, game_state);
      break;

    case "join":
      game_state.player_id = server_message.data.id;
      break;

    case "game_update":
      handleGameUpdate(server_message, game_state);
      break;

    case "game_end":
      handleGameEnd(server_message, game_state);

      break;
  }
}
function handleGameEnd(server_message, game_state) {
  const isDraw = server_message.data.winner === "draw";
  const winner =
    server_message.data.winner === game_state.player_id ? "player" : "opponent";

  game_container.setGameEnd(
    isDraw ? "draw" : winner,
    server_message.data.coords
  );
}

function handleGameUpdate(server_message, game_state) {
  const { x, y, state, active_player } = server_message.data;

  game_state.active_player = active_player;

  game_container.handleGameUpdate(
    { x, y },
    state,
    active_player === game_state.player_id
  );
}

function initializeGame(server_message, game_state) {
  game_container.setAttribute(
    "data-player-turn",
    server_message.data.starting_player === game_state.player_id
  );

  game_state.active_player = server_message.data.starting_player;
}

window.Auth.onAuthChange((user) => {
  const play_online_button = document.getElementById("play-online-button");

  if (user) {
    play_online_button.onclick = () => {
      openSocket();
    };
  } else {
    play_online_button.onclick = () => {
      window.Auth.login();
    };
  }
});
