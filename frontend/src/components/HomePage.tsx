import {JSX, useEffect, useState, useRef} from 'react'
import {ApiResponse, Course, Profile} from '../types/models'
import ProfileCreation from './ProfileCreation'
import CourseScanner from './CourseScanner'
import CourseDetail from './CourseDetail'
import Dashboard from './Dashboard'
import Settings from './Settings'
import './HomePage.css'

function HomePage(): JSX.Element {
    const [selectedProfile, setSelectedProfile] = useState<Profile | null>(null)
    const [profiles, setProfiles] = useState<Profile[]>([])
    const [courses, setCourses] = useState<Course[]>([])
    const [loading, setLoading] = useState<boolean>(true)
    const [showProfileCreation, setShowProfileCreation] = useState<boolean>(false)
    const [showCourseScanner, setShowCourseScanner] = useState<boolean>(false)
    const [selectedCourse, setSelectedCourse] = useState<string | null>(null)
    const [initialContentId, setInitialContentId] = useState<string | null>(null)
    const [profileStats, setProfileStats] = useState<{
        experience: number
        gems: number
        level: number
        streak: number
    } | null>(null)
    const [currentView, setCurrentView] = useState<'dashboard' | 'library' | 'settings'>('dashboard')

    const isFetchingStats = useRef(false)

    const baseURL = import.meta.env.VITE_BASE_URL

    // Fetch profiles on mount
    useEffect(() => {
        fetchProfiles()
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    // Fetch courses and stats when profile is selected (only on profile ID change)
    useEffect(() => {
        if (selectedProfile?.id) {
            fetchCourses()
            fetchProfileStats()
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [selectedProfile?.id])

    const fetchProfiles = async (): Promise<void> => {
        try {
            const response = await fetch(`${baseURL}/api/profiles`)
            const data: ApiResponse<Profile[]> = await response.json()
            if (data.success && data.data) {
                setProfiles(data.data)
            }
        } catch (error) {
            console.error('Error fetching profiles:', error)
        } finally {
            setLoading(false)
        }
    }

    const fetchCourses = async (): Promise<void> => {
        try {
            const response = await fetch(`${baseURL}/api/courses`)
            const data: ApiResponse<Course[]> = await response.json()
            if (data.success && data.data) {
                setCourses(data.data)
            }
        } catch (error) {
            console.error('Error fetching courses:', error)
        }
    }

    const fetchProfileStats = async (): Promise<void> => {
        if (!selectedProfile || isFetchingStats.current) return

        isFetchingStats.current = true
        try {
            const response = await fetch(`${baseURL}/api/profiles/${selectedProfile.id}`)
            const data: ApiResponse<Profile> = await response.json()
            if (data.success && data.data) {
                // Update only the stats state with fresh data
                const newStats = {
                    experience: data.data.experience || 0,
                    gems: data.data.gems || 0,
                    level: data.data.level || 1,
                    streak: data.data.streak || 0
                }
                console.log('Fetched profile stats:', newStats)
                setProfileStats(newStats)
            }
        } catch (error) {
            console.error('Error fetching profile stats:', error)
        } finally {
            isFetchingStats.current = false
        }
    }

    const handleProfileSelect = async (profile: Profile): Promise<void> => {
        try {
            await fetch(`${baseURL}/api/profiles/${profile.id}/select`, {
                method: 'POST',
            })
            setSelectedProfile(profile)
            // Initialize stats from the selected profile immediately
            setProfileStats({
                experience: profile.experience || 0,
                gems: profile.gems || 0,
                level: profile.level || 1,
                streak: profile.streak || 0
            })
        } catch (error) {
            console.error('Error selecting profile:', error)
        }
    }

    const handleLogout = (): void => {
        setSelectedProfile(null)
        setCourses([])
        setCurrentView('dashboard')
    }

    const handleProfileCreated = (): void => {
        setShowProfileCreation(false)
        fetchProfiles()
    }

    const handleCancelCreation = (): void => {
        setShowProfileCreation(false)
    }

    const handleCoursesImported = (): void => {
        setShowCourseScanner(false)
        fetchCourses()
    }

    const handleCancelScanner = (): void => {
        setShowCourseScanner(false)
    }

    const handleCourseClick = (courseId: string, contentId?: string): void => {
        setSelectedCourse(courseId)
        setInitialContentId(contentId ?? null)
    }

    const handleBackToCourses = (): void => {
        setSelectedCourse(null)
        setInitialContentId(null)
        fetchCourses()
        fetchProfileStats() // Refresh stats after completing content
        setCurrentView('dashboard')
    }

    const handleViewDashboard = (): void => {
        setCurrentView('dashboard')
    }

    const handleViewLibrary = (): void => {
        setCurrentView('library')
    }

    const handleViewSettings = (): void => {
        setCurrentView('settings')
    }

    const handleFactoryReset = (): void => {
        setSelectedProfile(null)
        setProfiles([])
        setCourses([])
        setProfileStats(null)
        setCurrentView('dashboard')
    }

    const handleDeleteProfile = async (profileId: string, profileName: string, e: React.MouseEvent): Promise<void> => {
        // Stop propagation to prevent profile selection
        e.stopPropagation()

        // Confirm deletion
        const confirmDelete = window.confirm(
            `Are you sure you want to delete the profile "${profileName}"?\n\nThis will permanently delete all progress and cannot be undone.`
        )

        if (!confirmDelete) {
            return
        }

        try {
            const response = await fetch(`${baseURL}/api/profiles`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({user_id: profileId}),
            })

            const data: ApiResponse<null> = await response.json()

            if (data.success) {
                // If deleted profile was selected, clear selection
                if (selectedProfile?.id === profileId) {
                    setSelectedProfile(null)
                    setCourses([])
                }
                // Refresh profile list
                fetchProfiles()
            } else {
                alert(`Failed to delete profile: ${data.message}`)
            }
        } catch (error) {
            console.error('Error deleting profile:', error)
            alert('Network error. Please try again.')
        }
    }

    if (loading) {
        return <div className="loading">Loading...</div>
    }

    // Show course detail view if a course is selected
    if (selectedCourse && selectedProfile) {
        return (
            <CourseDetail
                courseId={selectedCourse}
                userId={selectedProfile.id}
                onBack={handleBackToCourses}
                initialContentId={initialContentId ?? undefined}
            />
        )
    }

    return (
        <div className="homepage">
            {/* Navigation Bar */}
            <nav className="navbar">
                <div className="nav-brand">Course Management System</div>
                <div className="nav-links">
                    {selectedProfile && (
                        <>
                            <button
                                className={`nav-button ${currentView === 'dashboard' ? 'active' : ''}`}
                                onClick={handleViewDashboard}
                            >
                                🏠 Dashboard
                            </button>
                            <button
                                className={`nav-button ${currentView === 'library' ? 'active' : ''}`}
                                onClick={handleViewLibrary}
                            >
                                📚 Library
                            </button>
                        </>
                    )}
                    <button
                        className={`nav-button ${currentView === 'settings' ? 'active' : ''}`}
                        onClick={handleViewSettings}
                    >
                        ⚙️ Settings
                    </button>
                    {selectedProfile && (
                        <button className="nav-button logout" onClick={handleLogout}>
                            Switch Profile
                        </button>
                    )}
                </div>
            </nav>

            <main className="main-content">
                {currentView === 'settings' ? (
                    <Settings onFactoryReset={handleFactoryReset}/>
                ) : !selectedProfile ? (
                    // Profile Selection View
                    <div className="profile-selection">
                        <h1>Select Your Profile</h1>
                        <div className="profiles-grid">
                            {profiles.length > 0 ? (
                                profiles.map((profile) => (
                                    <div key={profile.id} className="profile-card">
                                        <button
                                            className="profile-button"
                                            onClick={() => handleProfileSelect(profile)}
                                        >
                                            <div className="profile-avatar">
                                                {profile.name.charAt(0).toUpperCase()}
                                            </div>
                                            <div className="profile-name">{profile.name}</div>
                                            <div className="profile-stats">
                                                <span>Level {profile.level || 1} </span>
                                                <span>⭐ {profile.experience || 0} XP</span>
                                                <span>💎 {profile.gems || 0}</span>
                                                <span>🔥 {profile.streak || 0} day streak</span>
                                            </div>
                                        </button>
                                        <button
                                            className="delete-profile-button"
                                            onClick={(e) => handleDeleteProfile(profile.id, profile.name, e)}
                                            title="Delete profile"
                                        >
                                            Delete Profile
                                        </button>
                                    </div>
                                ))
                            ) : (
                                <div className="no-profiles">
                                    <p>No profiles found. Create one to get started!</p>
                                    <button
                                        className="create-profile-button"
                                        onClick={() => setShowProfileCreation(true)}
                                    >
                                        Create New Profile
                                    </button>
                                </div>
                            )}
                        </div>

                        {/* Add Profile Creation button for existing profiles too */}
                        {profiles.length > 0 && (
                            <button
                                className="add-profile-button"
                                onClick={() => setShowProfileCreation(true)}
                            >
                                + Add Another Profile
                            </button>
                        )}
                    </div>
                ) : currentView === 'dashboard' ? (
                    // Dashboard View
                    <Dashboard
                        profile={{
                            ...selectedProfile,
                            experience: profileStats?.experience ?? selectedProfile.experience,
                            level: profileStats?.level ?? selectedProfile.level,
                            gems: profileStats?.gems ?? selectedProfile.gems,
                            streak: profileStats?.streak ?? selectedProfile.streak
                        }}
                        onCourseClick={handleCourseClick}
                        onScanCourses={() => setShowCourseScanner(true)}
                    />
                ) : (
                    // Library View
                    <div className="courses-view">
                        <div className="courses-header">
                            <div>
                                <h1>All Courses</h1>
                                <div className="user-stats">
                                    <span className="stat-badge level">
                                        ⭐ Level {profileStats?.level ?? 1}
                                    </span>
                                    <span className="stat-badge xp">
                                        💫 {profileStats?.experience ?? 0} XP
                                    </span>
                                    <span className="stat-badge gems">
                                        💎 {profileStats?.gems ?? 0} Gems
                                    </span>
                                    <span className="stat-badge streak">
                                        🔥 {profileStats?.streak ?? 0} Day Streak
                                    </span>
                                </div>
                            </div>
                        </div>

                        <div className="courses-section-header">
                            <h2>Course Library</h2>
                            <div style={{display: 'flex', gap: '10px'}}>
                                <button
                                    className="import-courses-button"
                                    onClick={() => setShowCourseScanner(true)}
                                >
                                    📥 Import Courses
                                </button>
                            </div>
                        </div>

                        <div className="courses-grid">
                            {courses.length > 0 ? (
                                courses.map((course) => (
                                    <div key={course.id} className="course-card">
                                        <div className="course-icon">📚</div>
                                        <h3>{course.title}</h3>
                                        <p className="course-description">
                                            {course.description || 'No description available'}
                                        </p>
                                        <div className="course-footer">
                                            <span className="course-modules">
                                                {course.modules?.length || 0} modules
                                            </span>
                                            <button
                                                className="course-button"
                                                onClick={() => handleCourseClick(course.id)}
                                            >
                                                Open
                                            </button>
                                        </div>
                                    </div>
                                ))
                            ) : (
                                <div className="no-courses">
                                    <p>No courses found. Import courses to get started!</p>
                                    <button
                                        className="import-courses-button-large"
                                        onClick={() => setShowCourseScanner(true)}
                                    >
                                        📥 Import Courses
                                    </button>
                                </div>
                            )}
                        </div>
                    </div>
                )}
            </main>

            {/* Profile Creation Modal */}
            {showProfileCreation && (
                <ProfileCreation
                    onProfileCreated={handleProfileCreated}
                    onCancel={handleCancelCreation}
                />
            )}

            {/* Course Scanner Modal */}
            {showCourseScanner && (
                <CourseScanner
                    onCoursesImported={handleCoursesImported}
                    onCancel={handleCancelScanner}
                />
            )}
        </div>
    )
}

export default HomePage