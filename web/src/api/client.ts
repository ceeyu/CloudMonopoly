import axios, { AxiosInstance, AxiosError } from 'axios';
import {
  CreateGameRequest,
  CreateGameResponse,
  JoinGameRequest,
  JoinGameResponse,
  GameStateResponse,
  TurnRequest,
  TurnResponse,
  DecisionRequest,
  DecisionResponse,
  ErrorResponse,
  EventResponse,
} from './types';

// API 基礎路徑 - 支援環境變數設定
// 部署時需要在 .env.production 中設定 VITE_API_URL
const getApiBaseUrl = () => {
  // 如果有設定環境變數，使用它
  if (import.meta.env.VITE_API_URL) {
    return `${import.meta.env.VITE_API_URL}/api/v1`;
  }
  // 本地開發使用相對路徑（會被 vite proxy 處理）
  return '/api/v1';
};

const API_BASE_URL = getApiBaseUrl();

// 建立 axios 實例
const apiClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 錯誤處理
function handleApiError(error: AxiosError<ErrorResponse>): never {
  if (error.response?.data) {
    throw new Error(error.response.data.message || error.response.data.error);
  }
  throw new Error(error.message || '網路錯誤');
}

// API 函數

/**
 * 建立新遊戲
 * POST /games
 */
export async function createGame(request: CreateGameRequest): Promise<CreateGameResponse> {
  try {
    const response = await apiClient.post<CreateGameResponse>('/games', request);
    return response.data;
  } catch (error) {
    throw handleApiError(error as AxiosError<ErrorResponse>);
  }
}

/**
 * 取得遊戲狀態
 * GET /games/:id
 */
export async function getGameState(gameId: string): Promise<GameStateResponse> {
  try {
    const response = await apiClient.get<GameStateResponse>(`/games/${gameId}`);
    return response.data;
  } catch (error) {
    throw handleApiError(error as AxiosError<ErrorResponse>);
  }
}

/**
 * 加入遊戲
 * POST /games/:id/join
 */
export async function joinGame(gameId: string, request: JoinGameRequest): Promise<JoinGameResponse> {
  try {
    const response = await apiClient.post<JoinGameResponse>(`/games/${gameId}/join`, request);
    return response.data;
  } catch (error) {
    throw handleApiError(error as AxiosError<ErrorResponse>);
  }
}

/**
 * 開始遊戲
 * POST /games/:id/start
 */
export async function startGame(gameId: string): Promise<{ message: string; game_id: string }> {
  try {
    const response = await apiClient.post<{ message: string; game_id: string }>(`/games/${gameId}/start`);
    return response.data;
  } catch (error) {
    throw handleApiError(error as AxiosError<ErrorResponse>);
  }
}

/**
 * 執行回合動作
 * POST /games/:id/turn
 */
export async function executeTurn(gameId: string, request: TurnRequest): Promise<TurnResponse> {
  try {
    const response = await apiClient.post<TurnResponse>(`/games/${gameId}/turn`, request);
    return response.data;
  } catch (error) {
    throw handleApiError(error as AxiosError<ErrorResponse>);
  }
}

/**
 * 提交決策
 * POST /games/:id/decision
 */
export async function submitDecision(gameId: string, request: DecisionRequest): Promise<DecisionResponse> {
  try {
    const response = await apiClient.post<DecisionResponse>(`/games/${gameId}/decision`, request);
    return response.data;
  } catch (error) {
    throw handleApiError(error as AxiosError<ErrorResponse>);
  }
}

/**
 * 取得事件
 * GET /games/:id/event
 */
export async function getEvent(gameId: string, eventId?: string): Promise<EventResponse> {
  try {
    const url = eventId ? `/games/${gameId}/event/${eventId}` : `/games/${gameId}/event`;
    const response = await apiClient.get<EventResponse>(url);
    return response.data;
  } catch (error) {
    throw handleApiError(error as AxiosError<ErrorResponse>);
  }
}

/**
 * 取得隨機事件
 * GET /games/:id/event/random/:type
 */
export async function getRandomEvent(gameId: string, eventType: string): Promise<EventResponse> {
  try {
    const response = await apiClient.get<EventResponse>(`/games/${gameId}/event/random/${eventType}`);
    return response.data;
  } catch (error) {
    throw handleApiError(error as AxiosError<ErrorResponse>);
  }
}

/**
 * 儲存遊戲
 * POST /games/:id/save
 */
export async function saveGame(gameId: string): Promise<{ message: string; game_id: string }> {
  try {
    const response = await apiClient.post<{ message: string; game_id: string }>(`/games/${gameId}/save`);
    return response.data;
  } catch (error) {
    throw handleApiError(error as AxiosError<ErrorResponse>);
  }
}

/**
 * 載入遊戲
 * GET /games/:id/load
 */
export async function loadGame(gameId: string): Promise<GameStateResponse> {
  try {
    const response = await apiClient.get<GameStateResponse>(`/games/${gameId}/load`);
    return response.data;
  } catch (error) {
    throw handleApiError(error as AxiosError<ErrorResponse>);
  }
}

/**
 * 健康檢查
 * GET /health
 */
export async function healthCheck(): Promise<{ status: string; service: string }> {
  try {
    const response = await apiClient.get<{ status: string; service: string }>('/health');
    return response.data;
  } catch (error) {
    throw handleApiError(error as AxiosError<ErrorResponse>);
  }
}

export default apiClient;
