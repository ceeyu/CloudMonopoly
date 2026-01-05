// API Types - 對應後端 Go 結構

// 公司類型
export type CompanyType = 'startup' | 'traditional' | 'cloud_reseller' | 'cloud_native';

// 遊戲狀態
export type GameStatus = 'waiting' | 'in_progress' | 'finished';

// 格子類型 (對應後端 CellType)
export type CellType = 'normal' | 'opportunity' | 'fate' | 'challenge' | 'start' | 'bonus';

// 格子類型資訊
export interface CellTypeInfo {
  name: string;
  icon: string;
  color: string;
  description: string;
}

export const CELL_TYPES: Record<CellType, CellTypeInfo> = {
  start: {
    name: '起點',
    icon: '🏁',
    color: '#10b981',
    description: '遊戲起點，經過可獲得獎勵',
  },
  normal: {
    name: '一般',
    icon: '📍',
    color: '#6b7280',
    description: '一般格子，獲得基礎資源',
  },
  opportunity: {
    name: '機會',
    icon: '🌟',
    color: '#f59e0b',
    description: '機會事件，可能獲得擴張或合作機會',
  },
  fate: {
    name: '命運',
    icon: '🎭',
    color: '#8b5cf6',
    description: '命運事件，可能遇到好運或挑戰',
  },
  challenge: {
    name: '關卡',
    icon: '⚔️',
    color: '#ef4444',
    description: '技術挑戰，需要做出架構決策',
  },
  bonus: {
    name: '獎勵',
    icon: '🎁',
    color: '#06b6d4',
    description: '獎勵格，獲得額外資源',
  },
};

// 公司資料
export interface Company {
  id: string;
  name: string;
  type: CompanyType;
  capital: number;
  employees: number;
  is_international: boolean;
  product_cycle: string;
  tech_debt: number;
  security_level: number;
  cloud_adoption: number;
  infrastructure: string[];
}

// 玩家狀態
export interface PlayerState {
  player_id: string;
  player_name: string;
  company: Company | null;
  position: number;
}

// 遊戲配置
export interface GameConfig {
  max_players: number;
  board_type: string;
  difficulty_level: string;
}

// 遊戲狀態回應
export interface GameStateResponse {
  game_id: string;
  status: GameStatus;
  current_turn: number;
  current_player_id: string;
  players: PlayerState[];
  board_size: number;
}

// 建立遊戲請求
export interface CreateGameRequest {
  max_players: number;
  board_type?: string;
  difficulty_level?: string;
}

// 建立遊戲回應
export interface CreateGameResponse {
  game_id: string;
  status: GameStatus;
  config: GameConfig;
  message: string;
}

// 加入遊戲請求
export interface JoinGameRequest {
  player_id: string;
  player_name: string;
  company_type: CompanyType;
}

// 加入遊戲回應
export interface JoinGameResponse {
  message: string;
  game_id: string;
}

// 回合動作請求
export interface TurnRequest {
  player_id: string;
  action_type: string;
  payload?: unknown;
}

// 回合結果回應
export interface TurnResponse {
  dice_value: number;
  old_position: number;
  new_position: number;
  capital_change: number;
  employee_change: number;
  circuit_completed: boolean;
  decision_required: boolean;
  cell_type: CellType;
}

// 決策請求
export interface DecisionRequest {
  player_id: string;
  event_id: string;
  choice_id: number;
}

// 決策回應
export interface DecisionResponse {
  success: boolean;
  message: string;
  explanation: string;
  capital_change: number;
  employee_change: number;
  learning_points: string[];
  aws_best_practice: string;
}

// 事件類型
export type EventType = 'opportunity' | 'fate' | 'challenge' | 'security';

// 事件背景
export interface EventContext {
  scenario: string;
  business_impact: string;
  technical_needs: string[];
  constraints: string[];
}

// 選項需求
export interface ChoiceRequirements {
  min_capital: number;
  min_employees: number;
  min_security_level: number;
  required_infra: string[];
}

// 選項結果
export interface ChoiceOutcomes {
  capital_change: number;
  employee_change: number;
  security_change: number;
  cloud_adoption_change: number;
  success_rate: number;
  time_to_implement: number;
}

// 事件選項
export interface EventChoice {
  id: number;
  title: string;
  description: string;
  is_aws: boolean;
  aws_services: string[];
  on_prem_solution: string;
  requirements: ChoiceRequirements;
  outcomes: ChoiceOutcomes;
  architecture_diagram: string;
}

// 事件回應
export interface EventResponse {
  id: string;
  type: EventType;
  title: string;
  description: string;
  real_world_ref: string;
  context: EventContext;
  choices: EventChoice[];
  aws_topics: string[];
}

// 錯誤回應
export interface ErrorResponse {
  error: string;
  code: string;
  message: string;
}

// 公司類型資訊 (對應後端 CompanyDefaults)
export interface CompanyTypeInfo {
  name: string;
  description: string;
  initialStats: {
    capital: number;
    employees: number;
    securityLevel: number;
    cloudAdoption: number;
  };
  icon: string;
}

export const COMPANY_TYPES: Record<CompanyType, CompanyTypeInfo> = {
  startup: {
    name: '新創公司',
    description: '資本較少但雲端採用率高，適合快速迭代',
    initialStats: {
      capital: 500,
      employees: 10,
      securityLevel: 2,
      cloudAdoption: 30,
    },
    icon: '🚀',
  },
  traditional: {
    name: '傳產公司',
    description: '資本雄厚但雲端採用率低，需要數位轉型',
    initialStats: {
      capital: 5000,
      employees: 200,
      securityLevel: 3,
      cloudAdoption: 10,
    },
    icon: '🏭',
  },
  cloud_reseller: {
    name: '雲端代理商',
    description: '熟悉雲端服務，資安等級高',
    initialStats: {
      capital: 2000,
      employees: 50,
      securityLevel: 4,
      cloudAdoption: 80,
    },
    icon: '🤝',
  },
  cloud_native: {
    name: '雲端公司',
    description: '完全雲端化，技術領先但資本有限',
    initialStats: {
      capital: 1000,
      employees: 30,
      securityLevel: 4,
      cloudAdoption: 95,
    },
    icon: '☁️',
  },
};
