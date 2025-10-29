import React, { useState, FormEvent } from 'react'
import { ApiResponse, Profile } from '../types/models'
import './ProfileCreation.css'

interface ProfileCreationProps {
  onProfileCreated: () => void
  onCancel: () => void
}

function ProfileCreation({ onProfileCreated, onCancel }: ProfileCreationProps) {
  const [name, setName] = useState<string>('')
  const [loading, setLoading] = useState<boolean>(false)
  const [error, setError] = useState<string>('')

  const baseURL = 'http://localhost:8080'

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()

    // Validate name
    if (!name.trim()) {
      setError('Please enter a name')
      return
    }

    setLoading(true)
    setError('')

    try {
      const response = await fetch(`${baseURL}/api/profiles`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name: name.trim() }),
      })

      const data: ApiResponse<Profile> = await response.json()

      if (data.success) {
        // Success! Close form and refresh profile list
        onProfileCreated()
      } else {
        setError(data.message || 'Failed to create profile')
      }
    } catch (err) {
      setError('Network error. Please try again.')
      console.error('Error creating profile:', err)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="profile-creation-overlay">
      <div className="profile-creation-container">
        <h2>Create New Profile</h2>

        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="profile-name">Profile Name</label>
            <input
              id="profile-name"
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Enter your name"
              disabled={loading}
              autoFocus
            />
          </div>

          {error && (
            <div className="error-message">
              {error}
            </div>
          )}

          <div className="form-buttons">
            <button
              type="button"
              onClick={onCancel}
              disabled={loading}
              className="cancel-button"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={loading || !name.trim()}
              className="submit-button"
            >
              {loading ? 'Creating...' : 'Create Profile'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

export default ProfileCreation
