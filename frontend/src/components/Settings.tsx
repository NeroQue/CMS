import {JSX, useState} from 'react'
import './Settings.css'

const baseURL = import.meta.env.VITE_BASE_URL

interface SettingsProps {
    onFactoryReset: () => void
}

interface DbStats {
    profiles: number
    courses: number
}

function Settings({onFactoryReset}: SettingsProps): JSX.Element {
    const [showConfirm, setShowConfirm] = useState<boolean>(false)
    const [confirmText, setConfirmText] = useState<string>('')
    const [resetting, setResetting] = useState<boolean>(false)
    const [stats, setStats] = useState<DbStats | null>(null)
    const [loadingStats, setLoadingStats] = useState<boolean>(false)
    const [message, setMessage] = useState<{ text: string; type: 'success' | 'error' } | null>(null)

    const CONFIRM_PHRASE = 'RESET'

    const fetchStats = async (): Promise<void> => {
        setLoadingStats(true)
        try {
            const response = await fetch(`${baseURL}/api/admin/stats`)
            const data = await response.json()
            if (data.success && data.data) {
                setStats(data.data)
            }
        } catch (error) {
            console.error('Error fetching stats:', error)
        } finally {
            setLoadingStats(false)
        }
    }

    const handleShowConfirm = (): void => {
        fetchStats()
        setShowConfirm(true)
        setConfirmText('')
        setMessage(null)
    }

    const handleCancel = (): void => {
        setShowConfirm(false)
        setConfirmText('')
        setMessage(null)
    }

    const handleFactoryReset = async (): Promise<void> => {
        if (confirmText !== CONFIRM_PHRASE) return

        setResetting(true)
        setMessage(null)

        try {
            const response = await fetch(`${baseURL}/api/admin/factory-reset`, {
                method: 'POST',
            })
            const data = await response.json()

            if (data.success) {
                setMessage({text: 'Database has been reset. All data cleared.', type: 'success'})
                setShowConfirm(false)
                setConfirmText('')
                // Wait a moment so the user sees the message, then trigger the reset callback
                setTimeout(() => {
                    onFactoryReset()
                }, 1500)
            } else {
                setMessage({text: data.message || 'Factory reset failed.', type: 'error'})
            }
        } catch (error) {
            console.error('Error during factory reset:', error)
            setMessage({text: 'Network error. Please try again.', type: 'error'})
        } finally {
            setResetting(false)
        }
    }

    return (
        <div className="settings-page">
            <h1>Settings</h1>

            {message && (
                <div className={`settings-message ${message.type}`}>
                    {message.type === 'success' ? '✅' : '❌'} {message.text}
                </div>
            )}

            <div className="settings-section">
                <div className="settings-section-header">
                    <h2>⚠️ Danger Zone</h2>
                </div>

                <div className="danger-zone">
                    <div className="danger-item">
                        <div className="danger-info">
                            <h3>Factory Reset Database</h3>
                            <p>
                                This will permanently delete <strong>all</strong> data: profiles,
                                courses, progress, and sessions. This action cannot be undone.
                            </p>
                        </div>
                        <button
                            className="btn-danger"
                            onClick={handleShowConfirm}
                            disabled={resetting}
                        >
                            🗑️ Factory Reset
                        </button>
                    </div>
                </div>
            </div>

            {/* Confirmation Modal */}
            {showConfirm && (
                <div className="modal-overlay" onClick={handleCancel}>
                    <div className="modal-content" onClick={(e) => e.stopPropagation()}>
                        <div className="modal-header">
                            <h2>☢️ Factory Reset</h2>
                        </div>
                        <div className="modal-body">
                            <p className="modal-message">
                                This will <strong>permanently delete all data</strong> from the database.
                            </p>

                            {loadingStats ? (
                                <div className="stats-loading">Loading database info...</div>
                            ) : stats && (
                                <div className="reset-stats">
                                    <p className="reset-stats-label">Data that will be deleted:</p>
                                    <ul className="reset-stats-list">
                                        <li>👤 {stats.profiles} profile{stats.profiles !== 1 ? 's' : ''}</li>
                                        <li>📚 {stats.courses} course{stats.courses !== 1 ? 's' : ''}</li>
                                        <li>📊 All progress records</li>
                                        <li>🔑 All sessions</li>
                                    </ul>
                                </div>
                            )}

                            <div className="confirm-input-section">
                                <label htmlFor="confirm-input">
                                    Type <strong>{CONFIRM_PHRASE}</strong> to confirm:
                                </label>
                                <input
                                    id="confirm-input"
                                    type="text"
                                    value={confirmText}
                                    onChange={(e) => setConfirmText(e.target.value)}
                                    placeholder={CONFIRM_PHRASE}
                                    autoFocus
                                    disabled={resetting}
                                />
                            </div>
                        </div>
                        <div className="modal-footer">
                            <button className="btn-cancel" onClick={handleCancel} disabled={resetting}>
                                Cancel
                            </button>
                            <button
                                className="btn-delete"
                                onClick={handleFactoryReset}
                                disabled={confirmText !== CONFIRM_PHRASE || resetting}
                            >
                                {resetting ? 'Resetting...' : '🗑️ Nuke Everything'}
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    )
}

export default Settings

