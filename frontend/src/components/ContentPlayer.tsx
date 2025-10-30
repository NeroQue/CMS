import {useEffect, useRef, useState} from 'react'
import './ContentPlayer.css'
import {ApiResponse, CompleteContentRequest, ContentItem, SaveProgressRequest, UserProgress,} from '../types/models'

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

    const handleMarkComplete = async () => {
        const requestBody: CompleteContentRequest = {
            user_id: userId,
        }

        try {
            const response = await fetch(`${baseURL}/api/content/${content.id}/complete`, {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify(requestBody),
            })

            const data: ApiResponse<UserProgress> = await response.json()

            if (data.success && data.data) {
                setIsCompleted(true)
                setProgress(data.data)
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
        </div>
    )
}

export default ContentPlayer

