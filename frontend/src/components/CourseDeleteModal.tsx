import {JSX} from 'react'
import './CourseDeleteModal.css'

interface CourseDeleteModalProps {
    courseName: string
    missingPaths?: string[]
    onConfirm: () => void
    onCancel: () => void
}

function CourseDeleteModal({courseName, missingPaths, onConfirm, onCancel}: CourseDeleteModalProps): JSX.Element {
    return (
        <div className="modal-overlay" onClick={onCancel}>
            <div className="modal-content" onClick={(e) => e.stopPropagation()}>
                <div className="modal-header">
                    <h2>⚠️ Course Not Found</h2>
                </div>
                <div className="modal-body">
                    <p className="modal-message">
                        The course <strong>"{courseName}"</strong> no longer exists on disk.
                    </p>
                    {missingPaths && missingPaths.length > 0 && (
                        <div className="missing-paths">
                            <p className="missing-paths-label">Missing files/folders:</p>
                            <ul className="missing-paths-list">
                                {missingPaths.slice(0, 5).map((path, index) => (
                                    <li key={index}>{path}</li>
                                ))}
                                {missingPaths.length > 5 && (
                                    <li className="more-items">...and {missingPaths.length - 5} more</li>
                                )}
                            </ul>
                        </div>
                    )}
                    <p className="modal-question">
                        Would you like to remove this course from the database?
                    </p>
                </div>
                <div className="modal-footer">
                    <button className="btn-cancel" onClick={onCancel}>
                        Keep Course
                    </button>
                    <button className="btn-delete" onClick={onConfirm}>
                        Delete from Database
                    </button>
                </div>
            </div>
        </div>
    )
}

export default CourseDeleteModal

