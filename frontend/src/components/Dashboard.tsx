import {JSX, useEffect, useState, useRef, useMemo} from 'react'
import './Dashboard.css'
import {Profile, Course, CourseProgress, ContentItem, ApiResponse} from '../types/models'

export const baseURL = import.meta.env.VITE_BASE_URL

interface DashboardProps {
    profile: Profile
    onCourseClick: (courseId: string, initialContentId?: string) => void
    onScanCourses?: () => void
}

interface LastAccessedContent {
    course: Course
    content: ContentItem
    progressPercent: number
}

function Dashboard({profile, onCourseClick, onScanCourses}: DashboardProps): JSX.Element {
    const [courses, setCourses] = useState<Course[]>([])
    const [progressMap, setProgressMap] = useState<Map<string, number>>(new Map())
    const [lastAccessedContent, setLastAccessedContent] = useState<LastAccessedContent | null>(null)
    const [loading, setLoading] = useState<boolean>(true)
    const [error, setError] = useState<string>('')

    const isFetching = useRef(false)

    // Fetch dashboard data
    useEffect(() => {
        isFetching.current = false

        fetchDashboardData().catch((err) => {
            console.error('Failed to fetch dashboard data:', err)
        })
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [profile.id])

    const fetchDashboardData = async () => {
        if (isFetching.current) return

        isFetching.current = true
        setLoading(true)
        setError('')

        try {
            // Fetch courses first
            const coursesRes = await fetch(`${baseURL}/api/courses`)
            const coursesData: ApiResponse<Course[]> = await coursesRes.json()

            if (!coursesData.success) {
                throw new Error('Failed to fetch courses')
            }

            // Treat empty DB result as an empty list
            const courseList = coursesData.data ?? []
            setCourses(courseList)

            // Fetch progress for each course in parallel
            const progressPromises = courseList.map(async (course) => {
                try {
                    const progressRes = await fetch(
                        `${baseURL}/api/courses/${course.id}/progress?user_id=${profile.id}`
                    )
                    const progressData: ApiResponse<CourseProgress> = await progressRes.json()

                    if (progressData.success && progressData.data) {
                        return {
                            courseId: course.id,
                            progress: progressData.data
                        }
                    }
                    return null
                } catch (err) {
                    console.error(`Failed to fetch progress for course ${course.id}:`, err)
                    return null
                }
            })

            const progressResults = await Promise.all(progressPromises)

            // Build progress map and find last accessed content
            const newProgressMap = new Map<string, number>()
            const validProgressResults: Array<{ courseId: string; progress: CourseProgress }> = []

            progressResults.forEach(result => {
                if (result && result.progress) {
                    newProgressMap.set(result.courseId, result.progress.completion_pct)
                    validProgressResults.push(result)
                }
            })

            setProgressMap(newProgressMap)

            // Find last accessed content from in-progress courses
            const inProgress = validProgressResults
                .filter(r => r.progress.completion_pct > 0 && r.progress.completion_pct < 100 && r.progress.last_accessed_at)
                .sort((a, b) => {
                    const dateA = a.progress.last_accessed_at ? new Date(a.progress.last_accessed_at).getTime() : 0
                    const dateB = b.progress.last_accessed_at ? new Date(b.progress.last_accessed_at).getTime() : 0
                    return dateB - dateA
                })

            if (inProgress.length > 0) {
                const mostRecent = inProgress[0]
                if (mostRecent) {
                    const recentCourse = courseList.find(c => c.id === mostRecent.courseId)

                    if (recentCourse?.modules && recentCourse.modules.length > 0) {
                        // Build a set of completed content item IDs from the progress data
                        const completedIds = new Set<string>()
                        if (mostRecent.progress.items) {
                            mostRecent.progress.items.forEach(item => {
                                if (item.completed) {
                                    completedIds.add(item.content_item_id)
                                }
                            })
                        }

                        // Walk modules in order to find the first incomplete content item
                        let foundContent: ContentItem | null = null
                        for (const module of recentCourse.modules) {
                            if (module.content_items && module.content_items.length > 0) {
                                for (const content of module.content_items) {
                                    if (!completedIds.has(content.id)) {
                                        foundContent = content
                                        break
                                    }
                                }
                                if (foundContent) break
                            }
                        }

                        // If all items are completed (shouldn't happen for in-progress), fall back to the first item
                        if (!foundContent) {
                            const firstModule = recentCourse.modules[0]
                            if (firstModule?.content_items && firstModule.content_items.length > 0) {
                                foundContent = firstModule.content_items[0] ?? null
                            }
                        }

                        if (foundContent) {
                            setLastAccessedContent({
                                course: recentCourse,
                                content: foundContent,
                                progressPercent: mostRecent.progress.completion_pct
                            })
                        }
                    }
                }
            }
        } catch (err) {
            setError('Failed to load dashboard data')
            console.error(err)
        } finally {
            setLoading(false)
            isFetching.current = false
        }
    }

    // Filter in-progress courses
    const inProgressCourses = useMemo(() => {
        return courses.filter(course => {
            const progress = progressMap.get(course.id) || 0
            return progress > 0 && progress < 100
        })
    }, [courses, progressMap])

    // Filter not started courses
    const notStartedCourses = useMemo(() => {
        return courses.filter(course => {
            const progress = progressMap.get(course.id) || 0
            return progress === 0
        })
    }, [courses, progressMap])

    const handleContinue = () => {
        if (lastAccessedContent) {
            onCourseClick(lastAccessedContent.course.id, lastAccessedContent.content.id)
        }
    }

    if (loading) {
        return (
            <div className="dashboard-skeleton">
                <div className="skeleton skeleton-header"></div>
                <div className="skeleton skeleton-card"></div>
                <div className="skeleton-grid">
                    {[1, 2, 3, 4].map(i => <div key={i} className="skeleton skeleton-course-card"></div>)}
                </div>
            </div>
        )
    }

    if (error) {
        return (
            <div className="dashboard-error">
                <p>{error}</p>
                <button onClick={() => {
                    fetchDashboardData().catch((err) => {
                        console.error('Failed to fetch dashboard data:', err)
                    })
                }}>Retry
                </button>
            </div>
        )
    }

    return (
        <div className="dashboard">
            {/* Stats Header */}
            <div className="stats-header">
                <div className="stat-card level">
                    <div className="stat-icon">⭐</div>
                    <div className="stat-content">
                        <div className="stat-value">{profile.level}</div>
                        <div className="stat-label">Level</div>
                    </div>
                </div>
                <div className="stat-card xp">
                    <div className="stat-icon">💫</div>
                    <div className="stat-content">
                        <div className="stat-value">{profile.experience}</div>
                        <div className="stat-label">XP</div>
                    </div>
                </div>
                <div className="stat-card gems">
                    <div className="stat-icon">💎</div>
                    <div className="stat-content">
                        <div className="stat-value">{profile.gems}</div>
                        <div className="stat-label">Gems</div>
                    </div>
                </div>
                <div className="stat-card streak">
                    <div className="stat-icon">🔥</div>
                    <div className="stat-content">
                        <div className="stat-value">{profile.streak}</div>
                        <div className="stat-label">Day Streak</div>
                    </div>
                </div>
            </div>

            {/* Continue Learning Card */}
            {lastAccessedContent && (
                <div className="continue-learning-card" onClick={handleContinue}>
                    <div className="card-header">
                        <h3>Continue Learning</h3>
                        <span className="time-info">Pick up where you left off</span>
                    </div>
                    <div className="card-body">
                        <p className="course-name">{lastAccessedContent.course.title}</p>
                        <p className="content-name">{lastAccessedContent.content.title}</p>
                        <div className="progress-bar">
                            <div className="progress-fill" style={{width: `${lastAccessedContent.progressPercent}%`}}/>
                        </div>
                        <span
                            className="progress-text">{lastAccessedContent.progressPercent.toFixed(0)}% complete</span>
                    </div>
                    <button className="resume-button">Resume ▶</button>
                </div>
            )}

            {/* In Progress Section */}
            {inProgressCourses.length > 0 && (
                <section className="in-progress-section">
                    <div className="section-header">
                        <h2>In Progress ({inProgressCourses.length})</h2>
                    </div>
                    <div className="course-grid">
                        {inProgressCourses.map(course => (
                            <div key={course.id} className="course-card" onClick={() => onCourseClick(course.id)}>
                                <div className="course-icon">📚</div>
                                <h3>{course.title}</h3>
                                <p className="course-description">{course.description || 'No description available'}</p>
                                <div className="course-progress">
                                    <div className="progress-bar">
                                        <div className="progress-fill"
                                             style={{width: `${progressMap.get(course.id) || 0}%`}}/>
                                    </div>
                                    <span className="progress-text">{(progressMap.get(course.id) || 0).toFixed(0)}% complete</span>
                                </div>
                                <div className="course-footer">
                                    <span className="course-modules">
                                        {course.modules?.length || 0} modules
                                    </span>
                                    <button className="course-button" onClick={(e) => {
                                        e.stopPropagation()
                                        onCourseClick(course.id)
                                    }}>Continue
                                    </button>
                                </div>
                            </div>
                        ))}
                    </div>
                </section>
            )}

            {/* All Courses / Not Started Section */}
            {notStartedCourses.length > 0 && (
                <section className="all-courses-section">
                    <div className="section-header">
                        <h2>Available Courses ({notStartedCourses.length})</h2>
                    </div>
                    <div className="course-grid">
                        {notStartedCourses.map(course => (
                            <div key={course.id} className="course-card" onClick={() => onCourseClick(course.id)}>
                                <div className="course-icon">📖</div>
                                <h3>{course.title}</h3>
                                <p className="course-description">{course.description || 'No description available'}</p>
                                <div className="course-footer">
                                    <span className="course-modules">
                                        {course.modules?.length || 0} modules
                                    </span>
                                    <button className="course-button" onClick={(e) => {
                                        e.stopPropagation()
                                        onCourseClick(course.id)
                                    }}>Start
                                    </button>
                                </div>
                            </div>
                        ))}
                    </div>
                </section>
            )}

            {/* Empty State */}
            {courses.length === 0 && (
                <div className="empty-state">
                    <div className="empty-icon">📚</div>
                    <h3 className="empty-title">No courses yet</h3>
                    <p className="empty-message">Import your first course to start learning!</p>
                    {onScanCourses && (
                        <button onClick={onScanCourses}>Scan for Courses</button>
                    )}
                </div>
            )}

            {/* Quick Actions */}
            {courses.length > 0 && onScanCourses && (
                <div className="quick-actions">
                    <button onClick={onScanCourses} className="quick-action-btn">
                        📥 Import More Courses
                    </button>
                </div>
            )}
        </div>
    )
}

export default Dashboard
