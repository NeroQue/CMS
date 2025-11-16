import {useEffect, useState} from 'react'
import './CourseDetail.css'
import {ApiResponse, ContentItem, Course, CourseProgress, Module,} from '../types/models'
import ContentPlayer from './ContentPlayer'

const baseURL = import.meta.env.VITE_BASE_URL

interface CourseDetailProps {
    courseId: string
    userId: string
    onBack: () => void
}

function CourseDetail({courseId, userId, onBack}: CourseDetailProps) {
    const [course, setCourse] = useState<Course | null>(null)
    const [progress, setProgress] = useState<CourseProgress | null>(null)
    const [selectedContent, setSelectedContent] = useState<ContentItem | null>(null)
    const [loading, setLoading] = useState<boolean>(true)
    const [error, setError] = useState<string>('')
    const [contentProgressMap, setContentProgressMap] = useState<Map<string, boolean>>(new Map())

    useEffect(() => {
        fetchCourseData()
    }, [courseId, userId])

    const fetchCourseData = async () => {
        setLoading(true)
        setError('')

        try {
            // Fetch course details with modules and content items
            const courseResponse = await fetch(`${baseURL}/api/courses/${courseId}`)
            const courseData: ApiResponse<Course> = await courseResponse.json()

            if (!courseData.success || !courseData.data) {
                setError('Failed to load course details')
                setLoading(false)
                return
            }

            setCourse(courseData.data)

            // Fetch course progress with batch optimization
            const progressResponse = await fetch(
                `${baseURL}/api/courses/${courseId}/progress?user_id=${userId}`
            )
            const progressData: ApiResponse<CourseProgress> = await progressResponse.json()

            if (progressData.success && progressData.data) {
                setProgress(progressData.data)

                // Build fast lookup map from batch progress data
                const completionMap = new Map<string, boolean>()

                if (progressData.data.items) {
                    progressData.data.items.forEach(item => {
                        completionMap.set(item.content_item_id, item.completed)
                    })
                }

                setContentProgressMap(completionMap)
            }

            setLoading(false)
        } catch (err) {
            console.error('Error fetching course data:', err)
            setError('Failed to load course. Please try again.')
            setLoading(false)
        }
    }

    const handleContentClick = (content: ContentItem) => {
        setSelectedContent(content)
    }

    const handleClosePlayer = () => {
        setSelectedContent(null)
        // Refresh progress after closing player
        fetchCourseData()
    }

    const getModuleProgress = (moduleId: string): number => {
        if (!progress || !progress.modules) return 0
        const moduleProgress = progress.modules.find((m) => m.module_id === moduleId)
        return moduleProgress ? moduleProgress.completion_pct : 0
    }

    const isContentCompleted = (contentId: string): boolean => {
        return contentProgressMap.get(contentId) === true
    }

    const handleResetProgress = async (): Promise<void> => {
        if (!window.confirm('Are you sure you want to reset all progress for this course? This cannot be undone.')) {
            return
        }

        try {
            const response = await fetch(
                `${baseURL}/api/courses/${courseId}/progress?user_id=${userId}`,
                {method: 'DELETE'}
            )
            const data: ApiResponse<any> = await response.json()

            if (data.success) {
                // Refresh course data to show updated progress
                fetchCourseData()
            } else {
                setError('Failed to reset progress')
            }
        } catch (err) {
            console.error('Error resetting progress:', err)
            setError('Failed to reset progress. Please try again.')
        }
    }

    if (loading) {
        return <div className="loading">Loading course...</div>
    }

    if (error) {
        return (
            <div className="error-container">
                <p className="error-message">{error}</p>
                <button onClick={onBack}>Go Back</button>
            </div>
        )
    }

    if (!course) {
        return (
            <div className="error-container">
                <p className="error-message">Course not found</p>
                <button onClick={onBack}>Go Back</button>
            </div>
        )
    }

    return (
        <div className="course-detail">
            {/* Header */}
            <div className="course-detail-header">
                <div className="header-top">
                    <button className="back-button" onClick={onBack}>
                        ← Back to Courses
                    </button>
                    <button className="reset-button" onClick={handleResetProgress}>
                        🔄 Reset Progress
                    </button>
                </div>
                <div className="course-header-info">
                    <h1>{course.title}</h1>
                    {course.description && <p className="course-desc">{course.description}</p>}
                    {progress && (
                        <div className="overall-progress">
                            <div className="progress-text">
                                Overall Progress: {progress.completed_items} / {progress.total_items} items
                                completed ({Math.round(progress.completion_pct)}%)
                            </div>
                            <div className="progress-bar">
                                <div
                                    className="progress-fill"
                                    style={{width: `${progress.completion_pct}%`}}
                                ></div>
                            </div>
                        </div>
                    )}
                </div>
            </div>

            {/* Modules and Content */}
            <div className="course-content">
                {!course.modules || course.modules.length === 0 ? (
                    <div className="no-modules">
                        <p>No modules found in this course.</p>
                    </div>
                ) : (
                    course.modules.map((module: Module) => (
                        <div key={module.id} className="module-section">
                            <div className="module-header">
                                <h2>{module.title}</h2>
                                {module.description && <p className="module-desc">{module.description}</p>}
                                <div className="module-progress">
                                    Progress: {Math.round(getModuleProgress(module.id))}%
                                </div>
                            </div>

                            <div className="content-list">
                                {!module.content_items || module.content_items.length === 0 ? (
                                    <div className="no-content">No content items in this module.</div>
                                ) : (
                                    module.content_items.map((content: ContentItem) => (
                                        <div
                                            key={content.id}
                                            className={`content-item ${
                                                isContentCompleted(content.id) ? 'completed' : ''
                                            }`}
                                            onClick={() => handleContentClick(content)}
                                        >
                                            <div className="content-icon">
                                                {content.content_type === 'video' && '🎥'}
                                                {content.content_type === 'pdf' && '📄'}
                                                {content.content_type === 'text' && '📝'}
                                                {!['video', 'pdf', 'text'].includes(content.content_type) && '📦'}
                                            </div>
                                            <div className="content-info">
                                                <div
                                                    className={`content-title ${isContentCompleted(content.id) ? 'completed' : ''}`}>{content.title}</div>
                                                {content.description && (
                                                    <div className="content-desc">{content.description}</div>
                                                )}
                                                <div className="content-meta">
                                                    {content.content_type}
                                                    {content.duration && ` • ${Math.ceil(content.duration / 60)} min`}
                                                </div>
                                            </div>
                                            {isContentCompleted(content.id) && (
                                                <div className="completed-badge">✓</div>
                                            )}
                                        </div>
                                    ))
                                )}
                            </div>
                        </div>
                    ))
                )}
            </div>

            {/* Content Player Modal */}
            {selectedContent && (
                <ContentPlayer
                    content={selectedContent}
                    userId={userId}
                    onClose={handleClosePlayer}
                />
            )}
        </div>
    )
}

export default CourseDetail
