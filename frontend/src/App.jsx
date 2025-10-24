import React, { useState } from 'react'
import './App.css'

function App() {
  const [output, setOutput] = useState('')
  const [realIDs, setRealIDs] = useState({
    userID: null,
    courseID: null,
    moduleID: null,
    contentID: null,
    taskID: null
  })

  const baseURL = 'http://localhost:8080'

  const testAPI = async (endpoint, method = 'GET', body = null, captureID = null) => {
    try {
      const options = {
        method,
        headers: {
          'Content-Type': 'application/json',
        },
      }

      if (body) {
        options.body = JSON.stringify(body)
      }

      const response = await fetch(`${baseURL}${endpoint}`, options)
      const data = await response.json()

      // Capture real IDs from successful responses
      if (captureID && data.success && data.data) {
        if (captureID === 'userID' && data.data.id) {
          setRealIDs(prev => ({ ...prev, userID: data.data.id }))
        } else if (captureID === 'courseID' && data.data.id) {
          setRealIDs(prev => ({ ...prev, courseID: data.data.id }))
          // Also capture module and content IDs if present in course data
          if (data.data.modules && data.data.modules.length > 0) {
            setRealIDs(prev => ({ ...prev, moduleID: data.data.modules[0].id }))
            if (data.data.modules[0].content_items && data.data.modules[0].content_items.length > 0) {
              setRealIDs(prev => ({ ...prev, contentID: data.data.modules[0].content_items[0].id }))
            }
          }
        } else if (captureID === 'profiles' && Array.isArray(data.data) && data.data.length > 0) {
          setRealIDs(prev => ({ ...prev, userID: data.data[0].id }))
        } else if (captureID === 'courses' && Array.isArray(data.data) && data.data.length > 0) {
          setRealIDs(prev => ({ ...prev, courseID: data.data[0].id }))
          // Try to capture nested module/content IDs from first course
          if (data.data[0].modules && data.data[0].modules.length > 0) {
            setRealIDs(prev => ({ ...prev, moduleID: data.data[0].modules[0].id }))
            if (data.data[0].modules[0].content_items && data.data[0].modules[0].content_items.length > 0) {
              setRealIDs(prev => ({ ...prev, contentID: data.data[0].modules[0].content_items[0].id }))
            }
          }
        } else if (captureID === 'taskID' && data.data.task_id) {
          setRealIDs(prev => ({ ...prev, taskID: data.data.task_id }))
        }
      }

      setOutput(`${method} ${endpoint}: ${JSON.stringify(data, null, 2)}`)
    } catch (error) {
      setOutput(`Error: ${error.message}`)
    }
  }

  const getWorkingUserID = () => realIDs.userID || 'no-user-created-yet'
  const getWorkingCourseID = () => realIDs.courseID || 'no-course-created-yet'
  const getWorkingModuleID = () => realIDs.moduleID || 'no-module-found-yet'
  const getWorkingContentID = () => realIDs.contentID || 'no-content-found-yet'
  const getWorkingTaskID = () => realIDs.taskID || 'no-task-created-yet'

  return (
    <div className="app">
      <header className="app-header">
        <h1>Course Management System - API Tester</h1>
        <div style={{ fontSize: '12px', color: '#666', marginTop: '10px' }}>
          Current User ID: {realIDs.userID || 'None (create a profile first)'}
          <br />
          Current Course ID: {realIDs.courseID || 'None (create a course first)'}
          <br />
          Current Module ID: {realIDs.moduleID || 'None (get courses with modules)'}
          <br />
          Current Content ID: {realIDs.contentID || 'None (get courses with content)'}
          <br />
          Current Task ID: {realIDs.taskID || 'None (start a batch import)'}
        </div>
      </header>
      <main>
        <div>
          <h2>API Endpoint Tests</h2>

          <h3>Basic</h3>
          <button onClick={() => testAPI('/api')}>Test Hello API</button>

          <h3>Profiles</h3>
          <button onClick={() => testAPI('/api/profiles', 'GET', null, 'profiles')}>Get Profiles (captures first profile ID)</button>
          <button onClick={() => testAPI('/api/profiles', 'POST', { name: 'Test User' }, 'userID')}>Create Profile (captures new ID)</button>
          <button onClick={() => testAPI('/api/profiles', 'PUT', { user_id: getWorkingUserID(), new_name: 'Updated User' })}>Update Profile (uses current ID)</button>
          <button onClick={() => testAPI('/api/profiles', 'DELETE', { user_id: getWorkingUserID() })}>Delete Profile (uses current ID)</button>
          <button onClick={() => testAPI(`/api/profiles/${getWorkingUserID()}/select`, 'POST')}>Select Profile (uses current ID)</button>

          <h3>Courses</h3>
          <button onClick={() => testAPI('/api/courses', 'GET', null, 'courses')}>Get Courses (captures IDs from modules/content)</button>
          <button onClick={() => testAPI('/api/courses', 'POST', { title: 'Test Course', description: 'A test course', relative_path: 'test-course' }, 'courseID')}>Create Course (captures new ID)</button>
          <button onClick={() => testAPI('/api/courses/directories')}>Get Course Directories</button>
          <button onClick={() => testAPI('/api/courses/scan')}>Scan New Courses</button>
          <button onClick={() => testAPI('/api/courses/batch', 'POST', { courses: [{ title: 'Batch Course', relative_path: 'batch-course' }] }, 'taskID')}>Batch Import (captures task ID)</button>

          <h3>Progress Tracking</h3>
          <button onClick={() => testAPI(`/api/courses/${getWorkingCourseID()}/progress?user_id=${getWorkingUserID()}`)}>Get Course Progress</button>
          <button onClick={() => testAPI(`/api/modules/${getWorkingModuleID()}/progress?user_id=${getWorkingUserID()}`)}>Get Module Progress (uses real module ID)</button>
          <button onClick={() => testAPI(`/api/content/${getWorkingContentID()}/progress`, 'POST', { user_id: getWorkingUserID(), progress_pct: 50.0 })}>Update Content Progress (uses real content ID)</button>
          <button onClick={() => testAPI(`/api/content/${getWorkingContentID()}/complete`, 'POST', { user_id: getWorkingUserID() })}>Mark Content Complete (uses real content ID)</button>
          <button onClick={() => testAPI(`/api/users/${getWorkingUserID()}/progress`)}>Get User Progress</button>

          <h3>Admin</h3>
          <button onClick={() => testAPI('/api/admin/factory-reset', 'POST')}>Factory Reset</button>
          <button onClick={() => testAPI('/api/admin/stats')}>Get Admin Stats</button>

          <h3>Tasks</h3>
          <button onClick={() => testAPI(`/api/tasks?id=${getWorkingTaskID()}`)}>Get Task Status (uses real task ID)</button>
          <button onClick={() => testAPI('/api/tasks/cleanup', 'POST')}>Cleanup Tasks</button>

          <h3>Output</h3>
          <pre style={{ background: '#3d3b3b', padding: '10px', maxHeight: '400px', overflow: 'auto', color: 'white' }}>
            {output || 'Click a button to test an endpoint...'}
          </pre>
        </div>
      </main>
    </div>
  )
}

export default App
