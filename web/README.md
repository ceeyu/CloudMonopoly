# AWS Learning Game - Web Frontend

AWS SAA 證照學習遊戲的 Web 前端應用。

## 技術棧

- React 18
- TypeScript
- Vite
- React Router
- Axios

## 開發

```bash
# 安裝依賴
npm install

# 啟動開發伺服器
npm run dev

# 建置生產版本
npm run build
```

## 專案結構

```
web/
├── src/
│   ├── api/           # API 客戶端和類型定義
│   │   ├── client.ts  # API 函數
│   │   └── types.ts   # TypeScript 類型
│   ├── components/    # React 元件
│   │   ├── GameBoard.tsx
│   │   ├── CompanyStatus.tsx
│   │   └── EventModal.tsx
│   ├── pages/         # 頁面元件
│   │   ├── Lobby.tsx  # 遊戲大廳
│   │   └── Game.tsx   # 遊戲主頁面
│   ├── App.tsx        # 主應用元件
│   └── main.tsx       # 入口點
├── public/            # 靜態資源
└── index.html         # HTML 模板
```

## API 連接

開發時，Vite 會將 `/api` 請求代理到 `http://localhost:8080`。

確保後端服務在 8080 端口運行：

```bash
# 在專案根目錄
go run cmd/server/main.go
```
