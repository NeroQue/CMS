import {useEffect} from 'react'
import './XPNotification.css'


interface XPNotificationProps {
    type: 'xp' | 'levelup' | 'gems'
    amount?: number
    oldLevel?: number
    newLevel?: number
    onClose: () => void
}

function XPNotification({type, amount, oldLevel, newLevel, onClose}: XPNotificationProps) {

    useEffect(() => {
        const timer = setTimeout(() => {
            onClose()
        }, 3000)

        return () => clearTimeout(timer)
    }, [onClose])

    const getIcon = (): string => {
        switch (type) {
            case 'xp':
                return '💫'
            case 'levelup':
                return '🎉'
            case 'gems':
                return '💎'
            default:
                return '✨'
        }
    }

    const getMessage = (): string => {
        switch (type) {
            case 'xp':
                return `+${amount} XP`
            case 'levelup':
                return `Level Up! Now Level ${newLevel}`
            case 'gems':
                return `+${amount} Gems`
            default:
                return 'Reward Earned!'
        }
    }

    const getClassName = (): string => {
        return `xp-notification xp-notification-${type}`
    }

    return (
        <div className={getClassName()}>
            <button className="notification-close" onClick={onClose}>×</button>
            <div className="notification-icon">{getIcon()}</div>
            <div className="notification-content">
                <div className="notification-message">{getMessage()}</div>
                {type === 'levelup' && oldLevel && (
                    <div className="notification-subtitle">
                        Level {oldLevel} → {newLevel}
                    </div>
                )}
            </div>
        </div>
    )
}

export default XPNotification


