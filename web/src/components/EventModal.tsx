import { useState } from 'react'
import { EventResponse, EventChoice, DecisionResponse } from '../api/types'
import './EventModal.css'

interface EventModalProps {
  event: EventResponse
  onDecision: (eventId: string, choiceId: number) => Promise<DecisionResponse>
  onClose: () => void
}

function EventModal({ event, onDecision, onClose }: EventModalProps) {
  const [selectedChoice, setSelectedChoice] = useState<number | null>(null)
  const [submitting, setSubmitting] = useState(false)
  const [decisionResult, setDecisionResult] = useState<DecisionResponse | null>(null)

  const handleSubmit = async () => {
    if (selectedChoice === null) return
    setSubmitting(true)
    try {
      const result = await onDecision(event.id, selectedChoice)
      setDecisionResult(result)
    } catch (error) {
      console.error('Decision failed:', error)
    } finally {
      setSubmitting(false)
    }
  }

  const handleCloseResult = () => {
    setDecisionResult(null)
    onClose()
  }

  // 顯示決策結果
  if (decisionResult) {
    return (
      <div className="modal-overlay" onClick={handleCloseResult}>
        <div className="modal-content result-modal" onClick={(e) => e.stopPropagation()}>
          <div className="modal-header">
            <h2>{decisionResult.success ? '✅ 決策成功' : '⚠️ 決策結果'}</h2>
            <button className="close-btn" onClick={handleCloseResult}>×</button>
          </div>

          <div className="modal-body">
            <div className="result-message">
              <p className="result-title">{decisionResult.message}</p>
              <p className="result-explanation">{decisionResult.explanation}</p>
            </div>

            <div className="result-changes">
              <h3>📊 影響</h3>
              <div className="changes-grid">
                <div className={`change-item ${decisionResult.capital_change >= 0 ? 'positive' : 'negative'}`}>
                  <span className="change-label">資本變化</span>
                  <span className="change-value">
                    {decisionResult.capital_change >= 0 ? '+' : ''}{decisionResult.capital_change} 萬
                  </span>
                </div>
                <div className={`change-item ${decisionResult.employee_change >= 0 ? 'positive' : 'negative'}`}>
                  <span className="change-label">員工變化</span>
                  <span className="change-value">
                    {decisionResult.employee_change >= 0 ? '+' : ''}{decisionResult.employee_change} 人
                  </span>
                </div>
              </div>
            </div>

            {decisionResult.learning_points && decisionResult.learning_points.length > 0 && (
              <div className="learning-section">
                <h3>📚 學習要點</h3>
                <ul className="learning-points">
                  {decisionResult.learning_points.map((point, index) => (
                    <li key={index}>{point}</li>
                  ))}
                </ul>
              </div>
            )}

            {decisionResult.aws_best_practice && (
              <div className="best-practice-section">
                <h3>💡 AWS 最佳實踐</h3>
                <p>{decisionResult.aws_best_practice}</p>
              </div>
            )}
          </div>

          <div className="modal-footer">
            <button className="btn btn-primary" onClick={handleCloseResult}>
              繼續遊戲
            </button>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2>📋 {event.title}</h2>
          <button className="close-btn" onClick={onClose}>×</button>
        </div>

        <div className="modal-body">
          <div className="event-description">
            <p>{event.description}</p>
            {event.real_world_ref && (
              <p className="real-world-ref">📖 參考案例: {event.real_world_ref}</p>
            )}
            <div className="event-context">
              <p><strong>情境：</strong>{event.context.scenario}</p>
              <p><strong>影響：</strong>{event.context.business_impact}</p>
              {event.context.technical_needs && event.context.technical_needs.length > 0 && (
                <p><strong>技術需求：</strong>{event.context.technical_needs.join('、')}</p>
              )}
              {event.context.constraints && event.context.constraints.length > 0 && (
                <p><strong>限制條件：</strong>{event.context.constraints.join('、')}</p>
              )}
            </div>
          </div>

          {event.aws_topics && event.aws_topics.length > 0 && (
            <div className="aws-topics">
              <span className="topics-label">相關 AWS 考點：</span>
              {event.aws_topics.map((topic, index) => (
                <span key={index} className="topic-tag">{topic}</span>
              ))}
            </div>
          )}

          <div className="choices-section">
            <h3>選擇方案</h3>
            <div className="choices-grid">
              {event.choices.map((choice: EventChoice) => (
                <div
                  key={choice.id}
                  className={`choice-card ${selectedChoice === choice.id ? 'selected' : ''}`}
                  onClick={() => setSelectedChoice(choice.id)}
                >
                  <div className="choice-header">
                    <h4>{choice.title}</h4>
                    {choice.is_aws && <span className="aws-badge">AWS</span>}
                  </div>
                  <p className="choice-desc">{choice.description}</p>
                  
                  {choice.aws_services && choice.aws_services.length > 0 && (
                    <div className="aws-services">
                      {choice.aws_services.map((service) => (
                        <span key={service} className="service-tag">{service}</span>
                      ))}
                    </div>
                  )}

                  {choice.on_prem_solution && !choice.is_aws && (
                    <p className="on-prem-solution">
                      <strong>地端方案：</strong>{choice.on_prem_solution}
                    </p>
                  )}

                  {choice.architecture_diagram && (
                    <div className="architecture-diagram">
                      <pre>{choice.architecture_diagram}</pre>
                    </div>
                  )}

                  <div className="choice-details">
                    <div className="requirements">
                      <strong>需求條件：</strong>
                      <ul>
                        {choice.requirements.min_capital > 0 && (
                          <li>最低資本: {choice.requirements.min_capital} 萬</li>
                        )}
                        {choice.requirements.min_employees > 0 && (
                          <li>最低員工: {choice.requirements.min_employees} 人</li>
                        )}
                        {choice.requirements.min_security_level > 0 && (
                          <li>資安等級: {choice.requirements.min_security_level}</li>
                        )}
                      </ul>
                    </div>
                    <div className="outcomes">
                      <strong>預期結果：</strong>
                      <ul>
                        <li className={choice.outcomes.capital_change >= 0 ? 'positive' : 'negative'}>
                          資本: {choice.outcomes.capital_change >= 0 ? '+' : ''}{choice.outcomes.capital_change} 萬
                        </li>
                        <li className={choice.outcomes.employee_change >= 0 ? 'positive' : 'negative'}>
                          員工: {choice.outcomes.employee_change >= 0 ? '+' : ''}{choice.outcomes.employee_change} 人
                        </li>
                        <li>成功率: {Math.round(choice.outcomes.success_rate * 100)}%</li>
                        <li>實施時間: {choice.outcomes.time_to_implement} 回合</li>
                      </ul>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        <div className="modal-footer">
          <button 
            className="btn btn-secondary" 
            onClick={onClose}
            disabled={submitting}
          >
            取消
          </button>
          <button 
            className="btn btn-primary"
            onClick={handleSubmit}
            disabled={selectedChoice === null || submitting}
          >
            {submitting ? '提交中...' : '確認決策'}
          </button>
        </div>
      </div>
    </div>
  )
}

export default EventModal
