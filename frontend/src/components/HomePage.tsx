import React, {useState, useEffect, JSX} from 'react'
import { Profile, Course, ApiResponse } from '../types/models'
import './HomePage.css'

function HomePage(): JSX.Element {
  const [selectedProfile, setSelectedProfile] = useState<Profile | null>(null)
  const [profiles, setProfiles] = useState<Profile[]>([])
  const [courses, setCourses] = useState<Course[]>([])
  const [loading, setLoading] = useState<boolean>(true)

  const baseURL = 'http://localhost:8080'

  // Fetch profiles on mount
  useEffect(() => {
    fetchProfiles()
  }, [])

  // Fetch courses when profile is selected
  useEffect(() => {
    if (selectedProfile) {
      fetchCourses()
    }
  }, [selectedProfile])

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

  const handleProfileSelect = async (profile: Profile): Promise<void> => {
    try {
      await fetch(`${baseURL}/api/profiles/${profile.id}/select`, {
        method: 'POST',
      })
      setSelectedProfile(profile)
    } catch (error) {
      console.error('Error selecting profile:', error)
    }
  }

  const handleLogout = (): void => {
    setSelectedProfile(null)
    setCourses([])
  }

  if (loading) {
    return <div className="loading">Loading...</div>
  }

  return (
    <div className="homepage">
      {/* Navigation Bar */}
      <nav className="navbar">
        <div className="nav-brand">Course Management System</div>
        <div className="nav-links">
          <button className="nav-button" disabled>Courses</button>
          <button className="nav-button" disabled>Progress</button>
          <button className="nav-button" disabled>Settings</button>
          {selectedProfile && (
            <button className="nav-button logout" onClick={handleLogout}>
              Switch Profile
            </button>
          )}
        </div>
      </nav>

      <main className="main-content">
        {!selectedProfile ? (
          // Profile Selection View
          <div className="profile-selection">
            <h1>Select Your Profile</h1>
            <div className="profiles-grid">
              {profiles.length > 0 ? (
                profiles.map((profile) => (
                  <button
                    key={profile.id}
                    className="profile-button"
                    onClick={() => handleProfileSelect(profile)}
                  >
                    <div className="profile-avatar">
                      {profile.name.charAt(0).toUpperCase()}
                    </div>
                    <div className="profile-name">{profile.name}</div>
                    <div className="profile-stats">
                      <span>⭐ {profile.experience || 0} XP</span>
                      <span>💎 {profile.gems || 0}</span>
                      <span>🔥 {profile.streak || 0} day streak</span>
                    </div>
                  </button>
                ))
              ) : (
                <div className="no-profiles">
                  <p>No profiles found. Create one to get started!</p>
                </div>
              )}
            </div>
          </div>
        ) : (
          // Courses View
          <div className="courses-view">
            <div className="courses-header">
              <div>
                <h1>Welcome back, {selectedProfile.name}!</h1>
                <div className="user-stats">
                  <span className="stat-badge">⭐ {selectedProfile.experience || 0} XP</span>
                  <span className="stat-badge">💎 {selectedProfile.gems || 0} Gems</span>
                  <span className="stat-badge">🔥 {selectedProfile.streak || 0} Day Streak</span>
                </div>
              </div>
            </div>

            <h2>Available Courses</h2>
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
                      <button className="course-button" disabled>
                        Start Course
                      </button>
                    </div>
                  </div>
                ))
              ) : (
                <div className="no-courses">
                  <p>No courses available yet. Check back soon!</p>
                </div>
              )}
            </div>
          </div>
        )}
      </main>
    </div>
  )
}

export default HomePage
