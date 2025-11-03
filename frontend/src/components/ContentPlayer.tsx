import {useEffect, useRef, useState} from 'react'
import './ContentPlayer.css'
import {
    ApiResponse,
    CompleteContentRequest,
    ContentItem,
    Profile,
    SaveProgressRequest,
    UserProgress,
} from '../types/models'
import XPNotification from './XPNotification'

const baseURL = 'http://localhost:8080'

interface ContentPlayerProps {
    content: ContentItem
    userId: string
    onClose: () => void
}

function ContentPlayer({content, userId, onClose}: ContentPlayerProps) {
    const [progress, setProgress] = useState<UserProgress | null>(null)
    const [loading, setLoading] = useState<boolean>(true)
    const [error, setError] = useState<string>('')
    const [isCompleted, setIsCompleted] = useState<boolean>(false)
    const videoRef = useRef<HTMLVideoElement>(null)
    const saveIntervalRef = useRef<NodeJS.Timeout | null>(null)
    const [notifications, setNotifications] = useState<Array<{
        id: string
        type: 'xp' | 'levelup' | 'gems'
        amount?: number
        oldLevel?: number
        newLevel?: number
    }>>([])
    const [, setProfileBeforeCompletion] = useState<Profile | null>(null)


    useEffect(() => {
        fetchProgress()
        return () => {
            // Clean up interval on unmount
            if (saveIntervalRef.current) {
                clearInterval(saveIntervalRef.current)
            }
        }
    }, [content.id, userId])

    useEffect(() => {
        // Set up auto-save for video progress
        if (content.content_type === 'video' && !isCompleted) {
            // Start auto-save interval every 10 seconds
            saveIntervalRef.current = setInterval(() => {
                if (videoRef.current && !videoRef.current.paused) {
                    saveVideoProgress()
                }
            }, 10000)
        }

        return () => {
            if (saveIntervalRef.current) {
                clearInterval(saveIntervalRef.current)
            }
        }
    }, [content.content_type, isCompleted])

    const fetchProgress = async () => {
        setLoading(true)
        setError('')

        try {
            const response = await fetch(
                `${baseURL}/api/content/${content.id}/progress?user_id=${userId}`
            )
            const data: ApiResponse<UserProgress> = await response.json()

            if (data.success && data.data) {
                setProgress(data.data)
                setIsCompleted(data.data.Completed)

                // If video, set starting position
                if (content.content_type === 'video' && videoRef.current && data.data.LastPosition?.Valid) {
                    videoRef.current.currentTime = data.data.LastPosition.Int32
                }
            }

            setLoading(false)
        } catch (err) {
            console.error('Error fetching progress:', err)
            setError('Failed to load progress')
            setLoading(false)
        }
    }

    const saveVideoProgress = async () => {
        if (!videoRef.current || content.content_type !== 'video') return

        const currentTime = videoRef.current.currentTime
        const duration = videoRef.current.duration || content.duration || 1
        const progressPct = Math.min((currentTime / duration) * 100, 100)

        const requestBody: SaveProgressRequest = {
            user_id: userId,
            last_position: Math.floor(currentTime),
            progress_pct: Math.floor(progressPct),
        }

        try {
            const response = await fetch(`${baseURL}/api/content/${content.id}/progress`, {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify(requestBody),
            })

            const data: ApiResponse<UserProgress> = await response.json()
            if (data.success && data.data) {
                setProgress(data.data)
            }
        } catch (err) {
            console.error('Error saving video progress:', err)
        }
    }

    const fetchProfileStats = async (): Promise<Profile | null> => {
        try {
            const response = await fetch(`${baseURL}/api/profiles/${userId}`)
            const data: ApiResponse<Profile> = await response.json()
            if (data.success && data.data) {
                return data.data
            }
            return null
        } catch (err) {
            console.error('Error fetching profile stats:', err)
            return null
        }
    }

    const addNotification = (notification: {
        type: 'xp' | 'levelup' | 'gems',
        amount?: number,
        oldLevel?: number,
        newLevel?: number
    }) => {
        const id = Date.now().toString() + Math.random().toString(36).substr(2, 9)
        setNotifications(prev => [...prev, {...notification, id}])
    }

    const removeNotification = (id: string) => {
        setNotifications(prev => prev.filter(n => n.id !== id))
    }

    const handleMarkComplete = async () => {
        // Fetch profile stats BEFORE marking complete
        const oldProfile = await fetchProfileStats()
        if (oldProfile) {
            setProfileBeforeCompletion(oldProfile)
        }

        const requestBody: CompleteContentRequest = {
            user_id: userId,
        }

        try {
            const response = await fetch(`${baseURL}/api/content/${content.id}/complete`, {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify(requestBody),
            })

            const data: ApiResponse<any> = await response.json()

            if (data.success && data.data) {
                setIsCompleted(true)
                setProgress(data.data.progress || data.data)

                // Extract XP data from response
                const xpAwarded = data.data.xp_awarded
                const xpAmount = data.data.xp_amount

                // Show XP notification if earned
                if (xpAwarded && xpAmount > 0) {
                    addNotification({type: 'xp', amount: xpAmount})
                }

                // Fetch new profile to check for level-up/gems
                const newProfile = await fetchProfileStats()
                if (oldProfile && newProfile) {
                    // Calculate levels
                    const oldLevel = Math.floor(oldProfile.experience / 100) + 1
                    const newLevel = Math.floor(newProfile.experience / 100) + 1

                    // Check for level up
                    if (newLevel > oldLevel) {
                        addNotification({
                            type: 'levelup',
                            oldLevel: oldLevel,
                            newLevel: newLevel
                        })
                    }

                    // Check for gems
                    const gemsEarned = newProfile.gems - oldProfile.gems
                    if (gemsEarned > 0) {
                        addNotification({type: 'gems', amount: gemsEarned})
                    }
                }
            } else {
                setError('Failed to mark as complete')
            }
        } catch (err) {
            console.error('Error marking complete:', err)
            setError('Failed to mark as complete')
        }
    }

    const handleVideoEnd = () => {
        // Auto-mark as complete when video ends
        if (!isCompleted) {
            handleMarkComplete()
        }
    }

    const getContentUrl = (): string => {
        // Construct URL to serve the file
        return `${baseURL}/api/content/${content.id}/file`
    }

    const renderContent = () => {
        if (loading) {
            return <div className="player-loading">Loading content...</div>
        }

        if (error) {
            return <div className="player-error">{error}</div>
        }

        switch (content.content_type) {
            case 'video':
                return (
                    <div className="video-container">
                        <video
                            ref={videoRef}
                            controls
                            className="video-player"
                            onEnded={handleVideoEnd}
                        >
                            <source src={getContentUrl()} type="video/mp4"/>
                            Your browser does not support the video tag.
                        </video>
                    </div>
                )

            case 'pdf':
                return (
                    <div className="pdf-container">
                        <iframe
                            src={getContentUrl()}
                            className="pdf-viewer"
                            title={content.title}
                        />
                    </div>
                )

            case 'text':
                return (
                    <div className="text-container">
                        <iframe
                            src={getContentUrl()}
                            className="text-viewer"
                            title={content.title}
                        />
                    </div>
                )

            default:
                return (
                    <div className="unsupported-content">
                        <p>This content type is not yet supported for viewing.</p>
                        <p>Content type: {content.content_type}</p>
                        <a href={getContentUrl()} download className="download-button">
                            Download File
                        </a>
                    </div>
                )
        }
    }

    return (
        <div className="content-player-overlay">
            <div className="content-player-modal">
                {/* Header */}
                <div className="player-header">
                    <div className="player-title-section">
                        <h2>{content.title}</h2>
                        {content.description && <p className="player-desc">{content.description}</p>}
                    </div>
                    <button className="close-button" onClick={onClose}>
                        ✕
                    </button>
                </div>

                {/* Content Area */}
                <div className="player-content">{renderContent()}</div>

                {/* Footer with completion and resume button */}
                <div className="player-footer">
                    {progress && progress.LastPosition?.Valid && progress.LastPosition.Int32 > 0 && !isCompleted && (
                        <button className="resume-button" onClick={() => {
                            if (videoRef.current && content.content_type === 'video') {
                                // @ts-ignore
                                videoRef.current.currentTime = progress.LastPosition.Int32;
                                videoRef.current.play();
                            }
                        }}>
                            Resume from {Math.floor(progress.LastPosition.Int32)}s
                        </button>
                    )}
                    {progress && (
                        <div className="progress-info">
                            Progress: {Math.round(progress.ProgressPct || 0)}%
                            {isCompleted && <span className="completed-text"> ✓ Completed</span>}
                        </div>
                    )}
                    {!isCompleted && (
                        <button className="complete-button" onClick={handleMarkComplete}>
                            Mark as Complete
                        </button>
                    )}
                    {isCompleted && (
                        <button className="completed-badge-button" disabled>
                            ✓ Completed
                        </button>
                    )}
                </div>
            </div>

            {/* XP Notifications */}
            {notifications.map(notification => (
                <XPNotification
                    key={notification.id}
                    type={notification.type}
                    amount={notification.amount}
                    oldLevel={notification.oldLevel}
                    newLevel={notification.newLevel}
                    onClose={() => removeNotification(notification.id)}
                />
            ))}
        </div>
    )
}

export default ContentPlayer

