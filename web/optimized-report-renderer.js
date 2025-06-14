/**
 * Optimized Playwright Report Renderer with Virtualization
 * 🚀 Performance improvements for large test datasets
 */

class VirtualizedTestRenderer {
    constructor(container, options = {}) {
        this.container = container;
        this.itemHeight = options.itemHeight || 60;
        this.bufferSize = options.bufferSize || 5;
        this.visibleItems = Math.ceil(window.innerHeight / this.itemHeight) + this.bufferSize * 2;
        
        this.allItems = [];
        this.filteredItems = [];
        this.scrollTop = 0;
        this.startIndex = 0;
        this.endIndex = 0;
        
        this.setupDOM();
        this.bindEvents();
    }
    
    setupDOM() {
        this.container.innerHTML = `
            <div class="virtualized-container" style="position: relative; overflow-y: auto; height: 100%;">
                <div class="spacer-before" style="height: 0px;"></div>
                <div class="visible-items" style="position: relative;"></div>
                <div class="spacer-after" style="height: 0px;"></div>
            </div>
        `;
        
        this.scrollContainer = this.container.querySelector('.virtualized-container');
        this.spacerBefore = this.container.querySelector('.spacer-before');
        this.visibleContainer = this.container.querySelector('.visible-items');
        this.spacerAfter = this.container.querySelector('.spacer-after');
    }
    
    bindEvents() {
        this.scrollContainer.addEventListener('scroll', 
            this.throttle(this.handleScroll.bind(this), 16));
        
        window.addEventListener('resize', 
            this.throttle(this.handleResize.bind(this), 250));
    }
    
    throttle(func, delay) {
        let timeoutId;
        let lastExecTime = 0;
        return function (...args) {
            const currentTime = Date.now();
            
            if (currentTime - lastExecTime > delay) {
                func.apply(this, args);
                lastExecTime = currentTime;
            } else {
                clearTimeout(timeoutId);
                timeoutId = setTimeout(() => {
                    func.apply(this, args);
                    lastExecTime = Date.now();
                }, delay - (currentTime - lastExecTime));
            }
        };
    }
    
    setData(items) {
        this.allItems = items;
        this.filteredItems = [...items];
        this.calculateIndices();
        this.render();
    }
    
    filter(predicate) {
        this.filteredItems = this.allItems.filter(predicate);
        this.scrollContainer.scrollTop = 0;
        this.calculateIndices();
        this.render();
    }
    
    handleScroll() {
        this.scrollTop = this.scrollContainer.scrollTop;
        this.calculateIndices();
        this.render();
    }
    
    handleResize() {
        this.visibleItems = Math.ceil(window.innerHeight / this.itemHeight) + this.bufferSize * 2;
        this.calculateIndices();
        this.render();
    }
    
    calculateIndices() {
        this.startIndex = Math.floor(this.scrollTop / this.itemHeight);
        this.startIndex = Math.max(0, this.startIndex - this.bufferSize);
        
        this.endIndex = this.startIndex + this.visibleItems;
        this.endIndex = Math.min(this.filteredItems.length, this.endIndex);
    }
    
    render() {
        const totalHeight = this.filteredItems.length * this.itemHeight;
        const offsetY = this.startIndex * this.itemHeight;
        
        this.spacerBefore.style.height = `${offsetY}px`;
        this.spacerAfter.style.height = `${totalHeight - offsetY - (this.endIndex - this.startIndex) * this.itemHeight}px`;
        
        // Clear visible container
        this.visibleContainer.innerHTML = '';
        
        // Render only visible items
        for (let i = this.startIndex; i < this.endIndex; i++) {
            const item = this.filteredItems[i];
            if (item) {
                const element = this.createItemElement(item, i);
                this.visibleContainer.appendChild(element);
            }
        }
    }
    
    createItemElement(testResult, index) {
        const div = document.createElement('div');
        div.className = `test-item ${testResult.status}`;
        div.style.height = `${this.itemHeight}px`;
        div.style.padding = '8px 16px';
        div.style.borderBottom = '1px solid #e0e0e0';
        div.style.display = 'flex';
        div.style.alignItems = 'center';
        div.style.cursor = 'pointer';
        
        const statusIcon = this.getStatusIcon(testResult.status);
        const duration = testResult.duration || 0;
        
        div.innerHTML = `
            <div style="flex: 0 0 24px; margin-right: 12px;">
                ${statusIcon}
            </div>
            <div style="flex: 1; min-width: 0;">
                <div style="font-weight: 500; color: #333; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
                    ${this.escapeHtml(testResult.title)}
                </div>
                <div style="font-size: 12px; color: #666; margin-top: 2px;">
                    ${testResult.file} • ${duration}ms
                </div>
            </div>
            <div style="flex: 0 0 auto; font-size: 12px; color: #666;">
                ${testResult.browser || 'chromium'}
            </div>
        `;
        
        div.addEventListener('click', () => this.onItemClick(testResult, index));
        
        return div;
    }
    
    getStatusIcon(status) {
        const icons = {
            passed: '<span style="color: #22c55e;">✓</span>',
            failed: '<span style="color: #ef4444;">✗</span>',
            skipped: '<span style="color: #f59e0b;">○</span>',
            timeout: '<span style="color: #8b5cf6;">⏱</span>'
        };
        return icons[status] || icons.passed;
    }
    
    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
    
    onItemClick(testResult, index) {
        // Lazy load test details
        this.loadTestDetails(testResult);
    }
    
    async loadTestDetails(testResult) {
        // Simulate loading test details
        const modal = document.createElement('div');
        modal.className = 'test-detail-modal';
        modal.style.cssText = `
            position: fixed; top: 0; left: 0; right: 0; bottom: 0;
            background: rgba(0,0,0,0.5); z-index: 1000;
            display: flex; align-items: center; justify-content: center;
        `;
        
        modal.innerHTML = `
            <div style="background: white; border-radius: 8px; padding: 24px; max-width: 80%; max-height: 80%; overflow-y: auto;">
                <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;">
                    <h3 style="margin: 0;">${this.escapeHtml(testResult.title)}</h3>
                    <button onclick="this.closest('.test-detail-modal').remove()" style="border: none; background: none; font-size: 24px; cursor: pointer;">×</button>
                </div>
                <div style="font-family: monospace; font-size: 14px; white-space: pre-wrap; background: #f5f5f5; padding: 16px; border-radius: 4px;">
                    ${this.escapeHtml(testResult.error || 'Test passed successfully')}
                </div>
            </div>
        `;
        
        document.body.appendChild(modal);
        
        modal.addEventListener('click', (e) => {
            if (e.target === modal) {
                modal.remove();
            }
        });
    }
}

// Paginated Report Manager
class PaginatedReportManager {
    constructor(container, options = {}) {
        this.container = container;
        this.pageSize = options.pageSize || 50;
        this.currentPage = 1;
        this.totalPages = 1;
        this.allData = [];
        this.filteredData = [];
        
        this.setupPagination();
    }
    
    setupPagination() {
        this.paginationContainer = document.createElement('div');
        this.paginationContainer.className = 'pagination-container';
        this.paginationContainer.style.cssText = `
            display: flex; justify-content: center; align-items: center;
            padding: 16px; gap: 8px; border-top: 1px solid #e0e0e0;
        `;
        
        this.container.appendChild(this.paginationContainer);
    }
    
    setData(data) {
        this.allData = data;
        this.filteredData = [...data];
        this.totalPages = Math.ceil(this.filteredData.length / this.pageSize);
        this.currentPage = 1;
        this.render();
    }
    
    filter(predicate) {
        this.filteredData = this.allData.filter(predicate);
        this.totalPages = Math.ceil(this.filteredData.length / this.pageSize);
        this.currentPage = 1;
        this.render();
    }
    
    getCurrentPageData() {
        const start = (this.currentPage - 1) * this.pageSize;
        const end = start + this.pageSize;
        return this.filteredData.slice(start, end);
    }
    
    goToPage(page) {
        if (page >= 1 && page <= this.totalPages) {
            this.currentPage = page;
            this.render();
        }
    }
    
    render() {
        this.renderPagination();
        this.renderCurrentPage();
    }
    
    renderPagination() {
        if (this.totalPages <= 1) {
            this.paginationContainer.style.display = 'none';
            return;
        }
        
        this.paginationContainer.style.display = 'flex';
        
        let paginationHTML = '';
        
        // Previous button
        paginationHTML += `
            <button onclick="window.reportManager.goToPage(${this.currentPage - 1})" 
                    ${this.currentPage === 1 ? 'disabled' : ''} 
                    style="padding: 8px 12px; border: 1px solid #ddd; background: white; cursor: pointer; border-radius: 4px;">
                Previous
            </button>
        `;
        
        // Page numbers
        const startPage = Math.max(1, this.currentPage - 2);
        const endPage = Math.min(this.totalPages, this.currentPage + 2);
        
        if (startPage > 1) {
            paginationHTML += `
                <button onclick="window.reportManager.goToPage(1)" 
                        style="padding: 8px 12px; border: 1px solid #ddd; background: white; cursor: pointer; border-radius: 4px;">
                    1
                </button>
            `;
            if (startPage > 2) {
                paginationHTML += `<span style="padding: 8px;">...</span>`;
            }
        }
        
        for (let i = startPage; i <= endPage; i++) {
            const isActive = i === this.currentPage;
            paginationHTML += `
                <button onclick="window.reportManager.goToPage(${i})" 
                        style="padding: 8px 12px; border: 1px solid #ddd; 
                               background: ${isActive ? '#007bff' : 'white'}; 
                               color: ${isActive ? 'white' : 'black'}; 
                               cursor: pointer; border-radius: 4px;">
                    ${i}
                </button>
            `;
        }
        
        if (endPage < this.totalPages) {
            if (endPage < this.totalPages - 1) {
                paginationHTML += `<span style="padding: 8px;">...</span>`;
            }
            paginationHTML += `
                <button onclick="window.reportManager.goToPage(${this.totalPages})" 
                        style="padding: 8px 12px; border: 1px solid #ddd; background: white; cursor: pointer; border-radius: 4px;">
                    ${this.totalPages}
                </button>
            `;
        }
        
        // Next button
        paginationHTML += `
            <button onclick="window.reportManager.goToPage(${this.currentPage + 1})" 
                    ${this.currentPage === this.totalPages ? 'disabled' : ''} 
                    style="padding: 8px 12px; border: 1px solid #ddd; background: white; cursor: pointer; border-radius: 4px;">
                Next
            </button>
        `;
        
        paginationHTML += `
            <span style="margin-left: 16px; color: #666; font-size: 14px;">
                Page ${this.currentPage} of ${this.totalPages} (${this.filteredData.length} items)
            </span>
        `;
        
        this.paginationContainer.innerHTML = paginationHTML;
    }
    
    renderCurrentPage() {
        const pageData = this.getCurrentPageData();
        // This would be called by the parent component to render the current page
        if (this.onPageRender) {
            this.onPageRender(pageData);
        }
    }
}

// Export for use
if (typeof module !== 'undefined' && module.exports) {
    module.exports = { VirtualizedTestRenderer, PaginatedReportManager };
} else {
    window.VirtualizedTestRenderer = VirtualizedTestRenderer;
    window.PaginatedReportManager = PaginatedReportManager;
}