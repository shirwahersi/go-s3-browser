// S3 Client Module - API wrapper
// Communicates with backend API instead of S3 directly

class S3Client {
    constructor() {
        // No configuration needed - backend handles credentials
    }

    // List objects at a given prefix (folder path)
    // Returns { folders: [], files: [] }
    async listObjects(prefix = '') {
        try {
            const response = await fetch('/api/list?prefix=' + encodeURIComponent(prefix));

            if (!response.ok) {
                throw new Error('Failed to list objects');
            }

            return await response.json();
        } catch (error) {
            console.error('Error listing objects:', error);
            throw error;
        }
    }

    // Get a pre-signed URL for downloading a file
    async getDownloadUrl(key) {
        try {
            const response = await fetch('/api/download?key=' + encodeURIComponent(key));

            if (!response.ok) {
                throw new Error('Failed to get download URL');
            }

            const data = await response.json();
            return data.url;
        } catch (error) {
            console.error('Error generating download URL:', error);
            throw error;
        }
    }

    // Helper: Format file size for display
    static formatSize(bytes) {
        if (bytes === 0) return '0 B';
        const units = ['B', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(1024));
        return parseFloat((bytes / Math.pow(1024, i)).toFixed(2)) + ' ' + units[i];
    }
}
