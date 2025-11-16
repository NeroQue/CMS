import {useEffect, useState} from 'react'
import {ApiResponse, CourseDirectory, CreateCourseInput, ScanResponse, Task} from '../types/models'
import './CourseScanner.css'

interface CourseScannerProps {
    onCoursesImported: () => void
    onCancel: () => void
}

function CourseScanner({onCoursesImported, onCancel}: CourseScannerProps) {
    const [scanning, setScanning] = useState<boolean>(false)
    const [directories, setDirectories] = useState<CourseDirectory[]>([])
    const [selectedDirs, setSelectedDirs] = useState<Set<string>>(new Set())
    const [importing, setImporting] = useState<boolean>(false)
    const [taskId, setTaskId] = useState<string | null>(null)
    const [taskStatus, setTaskStatus] = useState<Task | null>(null)
    const [error, setError] = useState<string>('')
    const [successMessage, setSuccessMessage] = useState<string>('')

    const baseURL = import.meta.env.VITE_BASE_URL

    // Cleanup interval on unmount
    useEffect(() => {
        let intervalId: NodeJS.Timeout | null = null

        if (taskId && importing) {
            intervalId = setInterval(() => {
                pollTaskStatus(taskId)
            }, 2000) // Poll every 2 seconds
        }

        return () => {
            if (intervalId) {
                clearInterval(intervalId)
            }
        }
    }, [taskId, importing])

    const handleScan = async () => {
        setScanning(true)
        setError('')
        setSuccessMessage('')
        setDirectories([])
        setSelectedDirs(new Set())

        try {
            const response = await fetch(`${baseURL}/api/courses/scan`)
            const data: ApiResponse<ScanResponse> = await response.json()

            console.log('Scan response:', data) // Debug log

            if (data.success && data.data) {
                const dirs = data.data.directories || []
                setDirectories(dirs)
                if (dirs.length === 0) {
                    setSuccessMessage('No new courses found. All courses are already imported!')
                } else {
                    setSuccessMessage(`Found ${dirs.length} new course(s) ready to import!`)
                }
            } else {
                setError(data.message || 'Failed to scan for courses')
            }
        } catch (err) {
            setError('Network error. Please check your connection.')
            console.error('Error scanning courses:', err)
        } finally {
            setScanning(false)
        }
    }

    const handleToggleSelect = (relativePath: string) => {
        const newSelected = new Set(selectedDirs)
        if (newSelected.has(relativePath)) {
            newSelected.delete(relativePath)
        } else {
            newSelected.add(relativePath)
        }
        setSelectedDirs(newSelected)
    }

    const handleSelectAll = () => {
        if (selectedDirs.size === directories.length) {
            setSelectedDirs(new Set())
        } else {
            setSelectedDirs(new Set(directories.map(d => d.relative_path)))
        }
    }

    const handleImport = async () => {
        if (selectedDirs.size === 0) {
            setError('Please select at least one course to import')
            return
        }

        setImporting(true)
        setError('')
        setSuccessMessage('')

        // Convert selected directories to CreateCourseInput format
        const coursesToImport: CreateCourseInput[] = Array.from(selectedDirs).map(relativePath => {
            const dir = directories.find(d => d.relative_path === relativePath)
            return {
                title: dir?.name || relativePath,
                relative_path: relativePath,
                description: `Imported from ${relativePath}`
            }
        })

        try {
            const response = await fetch(`${baseURL}/api/courses/batch`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({courses: coursesToImport}),
            })

            const data: ApiResponse<{ task_id: string }> = await response.json()

            if (data.success && data.data) {
                setTaskId(data.data.task_id)
                // Polling will start automatically via useEffect
            } else {
                setError(data.message || 'Failed to start import')
                setImporting(false)
            }
        } catch (err) {
            setError('Network error. Please try again.')
            console.error('Error starting import:', err)
            setImporting(false)
        }
    }

    const pollTaskStatus = async (id: string) => {
        try {
            const response = await fetch(`${baseURL}/api/tasks?id=${id}`)
            const data: ApiResponse<Task> = await response.json()

            if (data.success && data.data) {
                setTaskStatus(data.data)

                if (data.data.status === 'completed') {
                    handleImportComplete(data.data)
                } else if (data.data.status === 'failed') {
                    handleImportFailed(data.data)
                }
            }
        } catch (err) {
            console.error('Error polling task status:', err)
        }
    }

    const handleImportComplete = (task: Task) => {
        setImporting(false)
        setSuccessMessage(task.message || 'Courses imported successfully!')

        // Wait a bit then close and refresh
        setTimeout(() => {
            onCoursesImported()
        }, 2000)
    }

    const handleImportFailed = (task: Task) => {
        setImporting(false)
        setError(task.error_message || 'Import failed. Please try again.')
    }

    const getProgressPercentage = (): number => {
        if (!taskStatus) return 0
        return Math.round(taskStatus.progress || 0)
    }

    return (
        <div className="course-scanner-overlay">
            <div className="course-scanner-container">
                <div className="scanner-header">
                    <h2>Import New Courses</h2>
                    <button
                        className="close-button"
                        onClick={onCancel}
                        disabled={importing}
                    >
                        ×
                    </button>
                </div>

                <div className="scanner-content">
                    {/* Scan Button */}
                    {!importing && (
                        <button
                            className="scan-button"
                            onClick={handleScan}
                            disabled={scanning}
                        >
                            {scanning ? 'Scanning...' : 'Scan for New Courses'}
                        </button>
                    )}

                    {/* Error Message */}
                    {error && (
                        <div className="message error-message">
                            {error}
                        </div>
                    )}

                    {/* Success Message */}
                    {successMessage && !importing && (
                        <div className="message success-message">
                            {successMessage}
                        </div>
                    )}

                    {/* Directory List */}
                    {directories.length > 0 && !importing && (
                        <div className="directories-section">
                            <div className="directories-header">
                                <h3>Select Courses to Import</h3>
                                <button
                                    className="select-all-button"
                                    onClick={handleSelectAll}
                                >
                                    {selectedDirs.size === directories.length ? 'Deselect All' : 'Select All'}
                                </button>
                            </div>

                            <div className="directories-list">
                                {directories.map((dir) => (
                                    <label key={dir.relative_path} className="directory-item">
                                        <input
                                            type="checkbox"
                                            checked={selectedDirs.has(dir.relative_path)}
                                            onChange={() => handleToggleSelect(dir.relative_path)}
                                        />
                                        <div className="directory-info">
                                            <div className="directory-name">📚 {dir.name}</div>
                                            <div className="directory-path">{dir.relative_path}</div>
                                        </div>
                                    </label>
                                ))}
                            </div>

                            <button
                                className="import-button"
                                onClick={handleImport}
                                disabled={selectedDirs.size === 0}
                            >
                                Import {selectedDirs.size} Selected Course{selectedDirs.size !== 1 ? 's' : ''}
                            </button>
                        </div>
                    )}

                    {/* Import Progress */}
                    {importing && taskStatus && (
                        <div className="import-progress">
                            <h3>Importing Courses...</h3>
                            <div className="progress-bar-container">
                                <div
                                    className="progress-bar-fill"
                                    style={{width: `${getProgressPercentage()}%`}}
                                />
                            </div>
                            <div className="progress-info">
                                <span>{getProgressPercentage()}%</span>
                                <span>{taskStatus.message || 'Processing...'}</span>
                            </div>
                        </div>
                    )}
                </div>

                {/* Footer */}
                <div className="scanner-footer">
                    <button
                        className="cancel-footer-button"
                        onClick={onCancel}
                        disabled={importing}
                    >
                        {importing ? 'Close' : 'Cancel'}
                    </button>
                </div>
            </div>
        </div>
    )
}

export default CourseScanner
