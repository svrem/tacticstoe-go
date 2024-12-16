use wasm_bindgen::prelude::wasm_bindgen;

pub fn check_draw(board: [[i32; 4]; 4]) -> bool {
    for x in 0..4 {
        for y in 0..4 {
            if board[x][y] == 0 {
                return false;
            }
        }
    }
    true
}

#[derive(Debug)]
pub struct GameAction {
    pub player: i32,
    pub x: usize,
    pub y: usize,
}

pub fn check_for_move_validity(
    winner_data: &Option<BoardWinData>,
    action: &GameAction,
    board: [[i32; 4]; 4],
    next_board: [[i32; 4]; 4],
) -> bool {
    if board[action.x][action.y] != 0 {
        return false;
    }

    if winner_data.is_none() {
        return true;
    }

    let winner_data = winner_data.as_ref().unwrap();
    let opponent = if winner_data.player == 1 { 2 } else { 1 };

    // Check if the move is valid,
    // by checking if the opponent has placed their tick in a square directly adjacent
    // to the square the player is trying to place their tick in
    // e.g. if the player is trying to place their tick in (1, 1),
    // and the opponent has placed their tick in (0, 1), (2, 1), (1, 0), or (1, 2),
    // then the move is valid
    if (action.x == 0 || next_board[action.x - 1][action.y] != opponent)
        && (action.x == 3 || next_board[action.x + 1][action.y] != opponent)
        && (action.y == 0 || next_board[action.x][action.y - 1] != opponent)
        && (action.y == 3 || next_board[action.x][action.y + 1] != opponent)
    {
        return false;
    }

    true
}

pub struct BoardWinData {
    pub player: i32,
    coords: [[usize; 2]; 3],
}

pub fn check_board_for_win(board: [[i32; 4]; 4]) -> Option<BoardWinData> {
    for x in 0..4 {
        for y in 0..4 {
            if board[x][y] == 0 {
                continue;
            }

            if x > 0 && x < 3 && board[x][y] == board[x - 1][y] && board[x][y] == board[x + 1][y] {
                return Some(BoardWinData {
                    player: board[x][y],
                    coords: [[x - 1, y], [x, y], [x + 1, y]],
                });
            }

            if y > 0 && y < 3 && board[x][y] == board[x][y - 1] && board[x][y] == board[x][y + 1] {
                return Some(BoardWinData {
                    player: board[x][y],
                    coords: [[x, y - 1], [x, y], [x, y + 1]],
                });
            }

            if x > 0
                && x < 3
                && y > 0
                && y < 3
                && board[x][y] == board[x - 1][y - 1]
                && board[x][y] == board[x + 1][y + 1]
            {
                return Some(BoardWinData {
                    player: board[x][y],
                    coords: [[x - 1, y - 1], [x, y], [x + 1, y + 1]],
                });
            }

            if x > 0
                && x < 3
                && y > 0
                && y < 3
                && board[x][y] == board[x - 1][y + 1]
                && board[x][y] == board[x + 1][y - 1]
            {
                return Some(BoardWinData {
                    player: board[x][y],
                    coords: [[x - 1, y + 1], [x, y], [x + 1, y - 1]],
                });
            }
        }
    }

    None
}

type Board = [[i32; 4]; 4];

fn find_score_for_move(prev_board: Board, prev_action: GameAction, depth: i32, player: i32) -> i64 {
    if depth >= 5 {
        return 0;
    }

    let board = prev_board;

    let mut total_score: i64 = 0;
    for x in 0..4 {
        for y in 0..4 {
            if board[x][y] != 0 {
                continue;
            }

            let action = GameAction {
                x,
                y,
                player: if prev_action.player == 1 { 2 } else { 1 },
            };

            let mut next_board = board;
            next_board[action.x][action.y] = action.player;

            let winner_data = check_board_for_win(next_board);

            if !check_for_move_validity(&winner_data, &action, board, next_board) {
                continue;
            }

            if winner_data.is_some() {
                let winner_data = winner_data.unwrap();
                if winner_data.player == player {
                    total_score = 99999;
                } else {
                    total_score = -99999;
                }
            }

            total_score += find_score_for_move(next_board, action, depth + 1, player);
        }
    }

    if depth == 0 {
        return total_score;
    }

    return total_score / (depth.pow(4)) as i64;
}

#[wasm_bindgen]
pub fn find_best_move(board_str: String, player: i32) -> String {
    let mut board = [[0; 4]; 4];

    let mut i = 0;
    for x in 0..4 {
        for y in 0..4 {
            board[x][y] = board_str.chars().nth(i).unwrap().to_digit(10).unwrap() as i32;
            i += 1;
        }
    }

    let mut scores = vec![];

    for x in 0..4 {
        for y in 0..4 {
            if board[x][y] != 0 {
                scores.push(std::i64::MIN);
                continue;
            }

            let action = GameAction {
                x,
                y,
                player: player,
            };

            let mut next_board = board;
            next_board[action.x][action.y] = action.player;

            let score = find_score_for_move(next_board, action, 0, player);

            scores.push(score);
        }
    }

    let mut max_score = std::i64::MIN;
    let mut max_score_index = 0;

    for (i, score) in scores.iter().enumerate() {
        if *score > max_score {
            max_score = *score;
            max_score_index = i;
        }
    }

    let x = max_score_index / 4;
    let y = max_score_index % 4;

    return format!("{} {}", x, y);
}
