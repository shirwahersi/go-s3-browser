// Main Application Logic

class FileBrowser {
    constructor() {
        this.s3Client = null;
        this.currentPath = '';

        // DOM elements
        this.breadcrumbEl = document.getElementById('breadcrumb');
        this.fileListEl = document.getElementById('file-list');
        this.loadingEl = document.getElementById('loading');
        this.errorEl = document.getElementById('error');

        this.init();
    }

    async init() {
        try {
            this.s3Client = new S3Client();

            // Listen for browser back/forward navigation
            window.addEventListener('popstate', (event) => {
                const path = event.state?.path || '';
                this.navigateTo(path, false);
            });

            // Get initial path from URL or default to root
            const urlPath = new URLSearchParams(location.search).get('path') || '';

            // Set initial state (replaceState, not pushState)
            history.replaceState({ path: urlPath }, '', urlPath ? '?path=' + encodeURIComponent(urlPath) : location.pathname);

            await this.navigateTo(urlPath, false);
        } catch (error) {
            this.showError('Failed to initialize: ' + error.message);
        }
    }

    async navigateTo(path, pushState = true) {
        this.currentPath = path;
        this.showLoading(true);
        this.hideError();

        // Update browser history
        if (pushState) {
            const url = path ? '?path=' + encodeURIComponent(path) : location.pathname;
            history.pushState({ path }, '', url);
        }

        try {
            const { folders, files } = await this.s3Client.listObjects(path);
            this.renderBreadcrumb();
            this.renderFileList(folders, files);
        } catch (error) {
            this.showError('Failed to load folder: ' + error.message);
        } finally {
            this.showLoading(false);
        }
    }

    renderBreadcrumb() {
        const parts = this.currentPath.split('/').filter(p => p);

        let html = `<span class="breadcrumb-item" data-path="">Root</span>`;

        let accumulatedPath = '';
        for (const part of parts) {
            accumulatedPath += part + '/';
            html += `<span class="breadcrumb-separator">/</span>`;
            html += `<span class="breadcrumb-item" data-path="${accumulatedPath}">${part}</span>`;
        }

        this.breadcrumbEl.innerHTML = html;

        // Add click handlers
        this.breadcrumbEl.querySelectorAll('.breadcrumb-item').forEach(item => {
            item.addEventListener('click', () => {
                this.navigateTo(item.dataset.path);
            });
        });
    }

    renderFileList(folders, files) {
        if (folders.length === 0 && files.length === 0) {
            this.fileListEl.innerHTML = '<div class="empty-message">This folder is empty</div>';
            return;
        }

        let html = '';

        // Render folders first
        for (const folder of folders) {
            html += `
                <div class="file-item folder" data-key="${folder.key}">
                    <span class="file-icon">📁</span>
                    <span class="file-name">${this.escapeHtml(folder.name)}</span>
                </div>
            `;
        }

        // Render files
        for (const file of files) {
            const size = S3Client.formatSize(file.size);
            html += `
                <div class="file-item file" data-key="${file.key}">
                    <span class="file-icon">${this.getFileIcon(file.name)}</span>
                    <span class="file-name">${this.escapeHtml(file.name)}</span>
                    <span class="file-size">${size}</span>
                </div>
            `;
        }

        this.fileListEl.innerHTML = html;

        // Add click handlers for folders
        this.fileListEl.querySelectorAll('.folder').forEach(item => {
            item.addEventListener('click', () => {
                this.navigateTo(item.dataset.key);
            });
        });

        // Add click handlers for files
        this.fileListEl.querySelectorAll('.file').forEach(item => {
            item.addEventListener('click', () => {
                this.downloadFile(item.dataset.key);
            });
        });
    }

    async downloadFile(key) {
        try {
            const url = await this.s3Client.getDownloadUrl(key);
            window.open(url, '_blank');
        } catch (error) {
            this.showError('Failed to download file: ' + error.message);
        }
    }

    getFileIcon(filename) {
        const ext = filename.split('.').pop().toLowerCase();
        const icons = {
            // Images
            'jpg': '🖼️', 'jpeg': '🖼️', 'png': '🖼️', 'gif': '🖼️', 'svg': '🖼️', 'webp': '🖼️',
            // Documents
            'pdf': '📄', 'doc': '📄', 'docx': '📄', 'txt': '📄', 'rtf': '📄',
            // Spreadsheets
            'xls': '📊', 'xlsx': '📊', 'csv': '📊',
            // Archives
            'zip': '📦', 'rar': '📦', 'tar': '📦', 'gz': '📦', '7z': '📦',
            // Code
            'js': '📜', 'ts': '📜', 'py': '📜', 'html': '📜', 'css': '📜', 'json': '📜',
            // Video
            'mp4': '🎬', 'mkv': '🎬', 'avi': '🎬', 'mov': '🎬',
            // Audio
            'mp3': '🎵', 'wav': '🎵', 'flac': '🎵', 'ogg': '🎵'
        };
        return icons[ext] || '📄';
    }

    showLoading(show) {
        this.loadingEl.style.display = show ? 'block' : 'none';
        this.fileListEl.style.display = show ? 'none' : 'block';
    }

    showError(message) {
        this.errorEl.textContent = message;
        this.errorEl.style.display = 'block';
    }

    hideError() {
        this.errorEl.style.display = 'none';
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
}

// Initialize app when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    new FileBrowser();
});
