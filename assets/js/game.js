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
          <button onclick="this.parentElement.parentElement.hide()">Close</button>
        </div>
  `;
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
      <div class="player-info" id="opponent-info">
      </div>

      <div class="game-board" id="game-board">
      </div>

      <div class="player-info" id="player-info">
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
    if (!cell_index) return;

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

    for (const el of this.querySelectorAll("board-modal")) {
      el.remove();
    }

    for (const el of this.querySelectorAll(".player-info")) {
      el.innerHTML = "";
    }

    this.removeAttribute("data-loading");
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

  setGameEnd(winner, coords, new_elo_rating) {
    const delta_elo = new_elo_rating - (window?.Auth?.user?.elo_rating || 0);

    switch (winner) {
      case "draw":
        this.showModal(
          "Draw!",
          `No one won this game. Your ELO changed by ${delta_elo} points.`
        );
        this.setAttribute("data-draw", true);
        break;
      case "aborted":
        this.showModal(
          "Game Aborted!",
          `The game was aborted. Your ELO changed by ${delta_elo} points.`
        );

        break;
      case "player":
        this.showModal("You Won!", `Your ELO changed by ${delta_elo} points.`);
        this.setAttribute("data-player-winner", true);
        break;
      case "opponent":
        this.showModal("You Lost!", `Your ELO changed by ${delta_elo} points.`);
        this.setAttribute("data-player-winner", false);
        break;
    }

    if (window.Auth.user) window.Auth.user.elo_rating = new_elo_rating;

    if (coords.length === 0) return;

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
    game_container.setAttribute("data-loading", true);
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
      game_container.removeAttribute("data-loading");
      initializeGame(server_message, game_state);
      break;

    case "join":
      const player_info = document.getElementById("player-info");

      player_info.innerHTML = `
        <img src="${window.Auth.user.profile_picture}" alt="Player Profile Picture" onerror="window.onProfilePictureError(event)" />
        
        <div>
          <h3>${window.Auth.user.username}</h3>
          <p>${window.Auth.user.elo_rating}</p>
        </div>
      `;

      game_state.player_id = server_message.data.id;
      break;

    case "game_update":
      handleGameUpdate(server_message, game_state);
      break;

    case "game_end":
      handleGameEnd(server_message, game_state);

      break;

    default:
      console.error("Unknown message type: ", server_message.type);
  }
}
function handleGameEnd(server_message, game_state) {
  const new_elo_rating = server_message.data.new_elo_rating;

  const isDraw = server_message.data.winner === "draw";
  const isAbort = server_message.data.winner === "aborted";
  const winner =
    server_message.data.winner === game_state.player_id ? "player" : "opponent";

  game_container.setGameEnd(
    isDraw ? "draw" : isAbort ? "aborted" : winner,
    server_message.data.coords,
    new_elo_rating
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

  const opponent_info = document.querySelector("#opponent-info");

  opponent_info.innerHTML = `
    <img src="${server_message.data.opponent_picture}" alt="Opponent Profile Picture" onerror="window.onProfilePictureError(event)" />

    <div>
      <h3>${server_message.data.opponent_username}</h3>
      <p>${server_message.data.opponent_elo}</p>
    </div>
  `;

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

function checkForMoveValidity(action, board) {
  if (board[action.x][action.y] !== 0) {
    return false;
  }

  const nextBoard = board.map((row) => row.slice());
  nextBoard[action.x][action.y] = action.player;
  const winner_data = checkBoardForWin(nextBoard);

  if (!winner_data) return true;

  const opponent = winner_data.player === 1 ? 2 : 1;

  // Check if the move is valid,
  // by checking if the opponent has placed their tick in a square directly adjacent
  // to the square the player is trying to place their tick in
  // e.g. if the player is trying to place their tick in (1, 1),
  // and the opponent has placed their tick in (0, 1), (2, 1), (1, 0), or (1, 2),
  // then the move is valid
  if (
    (action.x === 0 || nextBoard[action.x - 1][action.y] !== opponent) &&
    (action.x === 3 || nextBoard[action.x + 1][action.y] !== opponent) &&
    (action.y === 0 || nextBoard[action.x][action.y - 1] !== opponent) &&
    (action.y === 3 || nextBoard[action.x][action.y + 1] !== opponent)
  )
    return false;

  return true;
}

function checkBoardForWin(board) {
  for (let x = 0; x < 4; x++) {
    for (let y = 0; y < 4; y++) {
      if (board[x][y] === 0) {
        continue;
      }

      if (
        x > 0 &&
        x < 3 &&
        board[x][y] === board[x - 1][y] &&
        board[x][y] === board[x + 1][y]
      ) {
        return {
          player: board[x][y],
          coords: [
            [x - 1, y],
            [x, y],
            [x + 1, y],
          ],
        };
      }

      if (
        y > 0 &&
        y < 3 &&
        board[x][y] === board[x][y - 1] &&
        board[x][y] === board[x][y + 1]
      ) {
        return {
          player: board[x][y],
          coords: [
            [x, y - 1],
            [x, y],
            [x, y + 1],
          ],
        };
      }

      if (
        x > 0 &&
        x < 3 &&
        y > 0 &&
        y < 3 &&
        board[x][y] === board[x - 1][y - 1] &&
        board[x][y] === board[x + 1][y + 1]
      ) {
        return {
          player: board[x][y],
          coords: [
            [x - 1, y - 1],
            [x, y],
            [x + 1, y + 1],
          ],
        };
      }

      if (
        x > 0 &&
        x < 3 &&
        y > 0 &&
        y < 3 &&
        board[x][y] === board[x - 1][y + 1] &&
        board[x][y] === board[x + 1][y - 1]
      ) {
        return {
          player: board[x][y],
          coords: [
            [x - 1, y + 1],
            [x, y],
            [x + 1, y - 1],
          ],
        };
      }
    }
  }

  return null;
}

const on_device_button = document.getElementById("play-on-device-button");

on_device_button.onclick = () => {
  game_container = document.getElementById("game-container");
  game_container.reset();

  game_container.setAttribute("data-player-turn", true);

  const board = Array.from({ length: 4 }, () => Array(4).fill(0));

  let active_player = 1;

  game_container.onCellClick(({ x, y }) => {
    const isMoveValid = checkForMoveValidity(
      { x, y, player: active_player },
      board
    );

    if (!isMoveValid) {
      return;
    }

    board[x][y] = active_player;
    game_container.handleGameUpdate({ x, y }, active_player, true);

    const winner_data = checkBoardForWin(board);

    if (winner_data) {
      game_container.setGameEnd(
        "player",
        winner_data.coords,
        window?.Auth?.user?.elo_rating || 0
      );
      return;
    }

    let isDraw = true;

    for (const row of board) {
      for (const cell of row) {
        if (cell === 0) {
          isDraw = false;
          break;
        }
      }
    }

    if (isDraw) {
      game_container.setGameEnd("draw", [], window.Auth?.user.elo_rating || 0);
      return;
    }

    active_player = active_player === 1 ? 2 : 1;
  });
};

// const play_bot_button = document.getElementById("play-bot-button");

// play_bot_button.onclick = async () => {
//   // import("/assets/pkg/bot_wasm.js").then((module) => {
//   //   module.find_best_move("0000000000000000", 1);
//   // });

//   const bot_wasm = await import("/assets/pkg/bot_wasm.js");
//   await bot_wasm.default();

//   console.log(bot_wasm);

//   const best_move = bot_wasm.find_best_move("0000000000000000", 1);

//   console.log(best_move);
// };

// const play_bot_button = document.getElementById("play-bot-button");

// play_bot_button.onclick = async () => {
//   const bot_wasm = await import("/assets/pkg/bot_wasm.js");
//   await bot_wasm.default();

//   game_container = document.getElementById("game-container");
//   game_container.reset();

//   function botMove() {
//     let state = "";

//     for (const cell of game_container.game_board.children) {
//       const s = cell.getAttribute("data-state");

//       if (!s) {
//         state += "0";
//         continue;
//       }

//       if (s == "1") {
//         state += "1";
//       } else {
//         state += "2";
//       }
//     }

//     console.log(state);
//     const bot_res_str = bot_wasm.find_best_move(state, 2);
//     const [x_s, y_s] = bot_res_str.split(" ");
//     const x = parseInt(x_s);
//     const y = parseInt(y_s);

//     console.log(x, y, bot_res_str, x_s, y_s);

//     game_container.handleGameUpdate({ x, y }, 2, true);
//   }

//   game_container.onCellClick(({ x, y }) => {
//     game_container.handleGameUpdate({ x, y }, 1, false);
//     botMove();
//   });

// };

function showTutorial() {
  const tutorial_template = document.getElementById("tutorial-template");
  const tutorial_dialog = document.createElement("app-dialog");

  tutorial_dialog.innerHTML = tutorial_template.innerHTML;

  document.body.appendChild(tutorial_dialog);

  setTimeout(() => {
    tutorial_dialog.show();
  });
}

setTimeout(() => {
  if (localStorage.getItem("tutorial-seen") !== "true") {
    showTutorial();
    localStorage.setItem("tutorial-seen", "true");
  }
}, 1000);

document.getElementById("tutorial-button").addEventListener("click", () => {
  showTutorial();
});
