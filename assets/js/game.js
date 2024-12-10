const game_container = document.getElementById("game-container");

document
  .getElementById("queue-button")
  .addEventListener("click", () => openSocket());

function openSocket() {
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

  function handleCellClick(e) {
    const cell_index = e.target.getAttribute("data-cell");

    const x = cell_index % 4;
    const y = Math.floor(cell_index / 4);

    socket.send(
      JSON.stringify({
        type: "action",
        data: {
          x,
          y,
        },
      })
    );
  }

  for (const cell of game_container.children) {
    cell.addEventListener("click", handleCellClick);
  }
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
  if (server_message.data.winner === "draw") {
    game_container.setAttribute("data-draw", true);
    return;
  }

  const clientIsWinner = server_message.data.winner === game_state.player_id;
  game_container.setAttribute("data-player-winner", clientIsWinner);

  const coords = server_message.data.coords;
  let coords_index = 0;

  const i = setInterval(() => {
    const [x, y] = coords[coords_index];
    const cell_index = y * 4 + x;
    const cell = game_container.children[cell_index];

    cell.setAttribute("data-winner", clientIsWinner);

    coords_index++;

    if (coords_index >= coords.length) {
      clearInterval(i);
    }
  }, 300);
}

function handleGameUpdate(server_message, game_state) {
  const { x, y, state, active_player } = server_message.data;

  game_container.setAttribute(
    "data-player-turn",
    active_player === game_state.player_id
  );
  game_state.active_player = active_player;

  const cell_index = y * 4 + x;

  game_state.active_player = active_player;
  const cell = game_container.children[cell_index];

  cell.setAttribute("data-state", state);

  if (state === 1) {
    cell.innerHTML =
      '<svg  xmlns="http://www.w3.org/2000/svg"  width="72"  height="72"  viewBox="0 0 24 24"  fill="none"  stroke="black"  stroke-width="2"  stroke-linecap="round"  stroke-linejoin="round"  class="icon icon-tabler icons-tabler-outline icon-tabler-x"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M18 6l-12 12" /><path d="M6 6l12 12" /></svg>';
  } else if (state === 2) {
    cell.innerHTML =
      '<svg  xmlns="http://www.w3.org/2000/svg"  width="72"  height="72"  viewBox="0 0 24 24"  fill="none"  stroke="black"  stroke-width="2"  stroke-linecap="round"  stroke-linejoin="round"  class="icon icon-tabler icons-tabler-outline icon-tabler-circle"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M12 12m-9 0a9 9 0 1 0 18 0a9 9 0 1 0 -18 0" /></svg>';
  }
}

function initializeGame(server_message, game_state) {
  game_container.setAttribute(
    "data-player-turn",
    server_message.data.starting_player === game_state.player_id
  );

  game_state.active_player = server_message.data.starting_player;
}
