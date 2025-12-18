let ws = null;
let reconnectInterval = null;

// API Examples for each service
const apiExamples = {
    'restful-inventory-manager': [
        {
            title: 'Get All Items',
            method: 'GET',
            url: '/items',
            description: 'Retrieve all inventory items'
        },
        {
            title: 'Create Item',
            method: 'POST',
            url: '/items',
            description: 'Add a new inventory item',
            body: {
                name: "Laptop",
                description: "High-performance laptop",
                price: 1500.99,
                quantity: 10
            }
        }
    ],
    'restful-task-manager': [
        {
            title: 'Get All Tasks',
            method: 'GET',
            url: '/tasks',
            description: 'Retrieve all tasks'
        },
        {
            title: 'Create Task',
            method: 'POST',
            url: '/tasks',
            description: 'Create a new task',
            body: {
                title: "Write documentation",
                description: "Document all APIs",
                due_date: "2025-12-31",
                priority: "high"
            }
        }
    ],
    'graphql-inventory-manager': [
        {
            title: 'Query All Items',
            method: 'POST',
            url: '/graphql',
            description: 'Get all inventory items with GraphQL',
            body: {
                query: "{ items { id name price quantity } }"
            }
        },
        {
            title: 'Add Item Mutation',
            method: 'POST',
            url: '/graphql',
            description: 'Add new item via GraphQL mutation',
            body: {
                query: 'mutation { addItem(name: "Headphones", description: "Wireless headphones", price: 99.99, quantity: 15) { id name } }'
            }
        }
    ],
    'grpc-inventory-manager': [
        {
            title: 'gRPC Service Info',
            method: 'INFO',
            url: 'localhost:8082',
            description: 'Use grpcurl or Postman to test',
            command: 'grpcurl -plaintext localhost:8082 list'
        }
    ],
    'grpc-user-registration': [
        {
            title: 'gRPC Service Info',
            method: 'INFO',
            url: 'localhost:8084',
            description: 'Use grpcurl or Postman to test',
            command: 'grpcurl -plaintext localhost:8084 list'
        }
    ],
    'websocket-live-chat': [
        {
            title: 'WebSocket Connection',
            method: 'WS',
            url: 'ws://localhost:8086/ws',
            description: 'Connect using WebSocket client',
            body: {
                type: "chat",
                username: "user123",
                message: "Hello everyone!"
            }
        }
    ]
};

// Initialize WebSocket connection
function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws`;

    ws = new WebSocket(wsUrl);

    ws.onopen = () => {
        console.log('WebSocket connected');
        updateConnectionStatus(true);
        if (reconnectInterval) {
            clearInterval(reconnectInterval);
            reconnectInterval = null;
        }
    };

    ws.onmessage = (event) => {
        const services = JSON.parse(event.data);
        renderServices(services);
    };

    ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        updateConnectionStatus(false);
    };

    ws.onclose = () => {
        console.log('WebSocket disconnected');
        updateConnectionStatus(false);

        // Attempt reconnection
        if (!reconnectInterval) {
            reconnectInterval = setInterval(() => {
                console.log('Attempting to reconnect...');
                connectWebSocket();
            }, 3000);
        }
    };
}

function updateConnectionStatus(connected) {
    const indicator = document.getElementById('connection-indicator');
    const text = document.getElementById('connection-text');

    if (connected) {
        indicator.classList.remove('disconnected');
        indicator.classList.add('connected');
        text.textContent = 'Connected';
    } else {
        indicator.classList.remove('connected');
        indicator.classList.add('disconnected');
        text.textContent = 'Disconnected';
    }
}

let selectedServiceId = null;
let allServices = [];

function renderServices(services) {
    allServices = services;
    const stoppedColumn = document.getElementById('stopped-services');
    const runningColumn = document.getElementById('running-services');
    const template = document.getElementById('service-card-template');

    // Clear existing cards
    stoppedColumn.innerHTML = '';
    runningColumn.innerHTML = '';

    // Separate services by status
    const stoppedServices = services.filter(s => s.status !== 'running');
    const runningServices = services.filter(s => s.status === 'running');

    // Render stopped services
    stoppedServices.forEach(service => {
        const card = createServiceCard(template, service);
        stoppedColumn.appendChild(card);
    });

    // Render running services
    runningServices.forEach(service => {
        const card = createServiceCard(template, service);
        runningColumn.appendChild(card);
    });

    // Update test panel if a service was previously selected
    if (selectedServiceId) {
        const selectedService = services.find(s => s.id === selectedServiceId);
        if (selectedService && selectedService.status === 'running') {
            updateTestPanel(selectedService);
            // Re-apply selected class
            const selectedCard = document.querySelector(`[data-service-id="${selectedServiceId}"]`);
            if (selectedCard) {
                selectedCard.classList.add('selected');
            }
        } else {
            // Service is no longer running, clear selection
            selectedServiceId = null;
            showNoSelectionMessage();
        }
    }
}

function createServiceCard(template, service) {
    const card = template.content.cloneNode(true);
    const cardElement = card.querySelector('.service-card');

    cardElement.setAttribute('data-service-id', service.id);
    cardElement.setAttribute('data-service-port', service.port);
    cardElement.classList.remove('running', 'error', 'selected');

    if (service.status === 'running') {
        cardElement.classList.add('running');
    } else if (service.status === 'error') {
        cardElement.classList.add('error');
    }

    if (selectedServiceId === service.id) {
        cardElement.classList.add('selected');
    }

    card.querySelector('.service-name').textContent = service.name;
    card.querySelector('.service-type-badge').textContent = service.type;
    card.querySelector('.service-description').textContent = service.description;
    card.querySelector('.port-number').textContent = service.port;

    const statusIndicator = card.querySelector('.status-indicator');
    const statusText = card.querySelector('.status-text');

    statusIndicator.className = `status-indicator ${service.status}`;
    statusText.textContent = service.status;

    // Update buttons based on status
    const startBtn = card.querySelector('.btn-start');
    const stopBtn = card.querySelector('.btn-stop');

    if (service.status === 'running') {
        startBtn.disabled = true;
        stopBtn.disabled = false;
    } else {
        startBtn.disabled = false;
        stopBtn.disabled = true;
    }

    // Show error if present
    if (service.error) {
        const errorDiv = card.querySelector('.service-error');
        errorDiv.style.display = 'block';
        errorDiv.querySelector('.error-text').textContent = service.error;
    }

    return card;
}

async function startService(button) {
    const card = button.closest('.service-card');
    const serviceId = card.getAttribute('data-service-id');

    button.disabled = true;
    button.textContent = 'Starting...';

    try {
        const response = await fetch('/api/services/start', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ service_id: serviceId })
        });

        if (!response.ok) {
            const error = await response.text();
            throw new Error(error);
        }

        // Status will be updated via WebSocket
    } catch (error) {
        console.error('Error starting service:', error);
        alert(`Failed to start service: ${error.message}`);
        button.disabled = false;
        button.textContent = 'Start';
    }
}

async function stopService(button) {
    const card = button.closest('.service-card');
    const serviceId = card.getAttribute('data-service-id');

    button.disabled = true;
    button.textContent = 'Stopping...';

    try {
        const response = await fetch('/api/services/stop', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ service_id: serviceId })
        });

        if (!response.ok) {
            const error = await response.text();
            throw new Error(error);
        }

        // Status will be updated via WebSocket
    } catch (error) {
        console.error('Error stopping service:', error);
        alert(`Failed to stop service: ${error.message}`);
        button.disabled = false;
        button.textContent = 'Stop';
    }
}

// Initial load
async function loadInitialServices() {
    try {
        const response = await fetch('/api/services');
        const services = await response.json();
        renderServices(services);
    } catch (error) {
        console.error('Error loading services:', error);
    }
}

// Select a service and show its API testing panel
function selectService(cardElement) {
    const serviceId = cardElement.getAttribute('data-service-id');
    const service = allServices.find(s => s.id === serviceId);

    // Only allow selection of running services
    if (!service || service.status !== 'running') {
        return;
    }

    // Update selected service
    selectedServiceId = serviceId;

    // Remove selected class from all cards
    document.querySelectorAll('.service-card').forEach(card => {
        card.classList.remove('selected');
    });

    // Add selected class to clicked card
    cardElement.classList.add('selected');

    // Update test panel
    updateTestPanel(service);
}

function updateTestPanel(service) {
    const testPanelContent = document.getElementById('test-panel-content');
    const examples = apiExamples[service.id] || [];

    const examplesHTML = examples.map((example, index) => `
        <div class="api-example" data-example-index="${index}">
            <div class="example-header">
                <span class="method-badge ${example.method.toLowerCase()}">${example.method}</span>
                <strong>${example.title}</strong>
            </div>
            <p class="example-description">${example.description}</p>
            <div class="example-details">
                <code class="example-url">http://localhost:${service.port}${example.url}</code>
                ${example.body ? `<pre class="example-body">${JSON.stringify(example.body, null, 2)}</pre>` : ''}
                ${example.command ? `<pre class="example-command">${example.command}</pre>` : ''}
            </div>
            ${example.method === 'GET' || example.method === 'POST' ? `
                <button class="btn btn-try" onclick="tryApiCall('${service.id}', ${index}, ${service.port})">Try It</button>
            ` : ''}
        </div>
    `).join('');

    testPanelContent.innerHTML = `
        <div class="test-panel-service-header">
            <span class="service-type-badge">${service.type}</span>
            <h3>${service.name}</h3>
        </div>
        ${examplesHTML}
    `;
}

function showNoSelectionMessage() {
    const testPanelContent = document.getElementById('test-panel-content');
    testPanelContent.innerHTML = `
        <div class="no-selection">
            <p>Click on a running service to test its API endpoints</p>
        </div>
    `;
}

// Try an API call
async function tryApiCall(serviceId, exampleIndex, port) {
    const example = apiExamples[serviceId][exampleIndex];
    const url = `http://localhost:${port}${example.url}`;

    // Find the example div in the test panel
    const exampleDiv = document.querySelector(`#test-panel-content .api-example[data-example-index="${exampleIndex}"]`);
    if (!exampleDiv) return;

    // Add loading indicator
    const tryBtn = exampleDiv.querySelector('.btn-try');
    const originalText = tryBtn.textContent;
    tryBtn.textContent = 'Loading...';
    tryBtn.disabled = true;

    // Remove previous result if exists
    const oldResult = exampleDiv.querySelector('.api-result');
    if (oldResult) oldResult.remove();

    try {
        const options = {
            method: example.method,
            headers: { 'Content-Type': 'application/json' }
        };

        if (example.body) {
            options.body = JSON.stringify(example.body);
        }

        const response = await fetch(url, options);
        const data = await response.json();

        // Display result
        const resultHTML = `
            <div class="api-result success">
                <strong>✓ Response (${response.status}):</strong>
                <pre>${JSON.stringify(data, null, 2)}</pre>
            </div>
        `;
        exampleDiv.insertAdjacentHTML('beforeend', resultHTML);

    } catch (error) {
        const resultHTML = `
            <div class="api-result error">
                <strong>✗ Error:</strong>
                <pre>${error.message}</pre>
            </div>
        `;
        exampleDiv.insertAdjacentHTML('beforeend', resultHTML);
    } finally {
        tryBtn.textContent = originalText;
        tryBtn.disabled = false;
    }
}

// Initialize on page load
document.addEventListener('DOMContentLoaded', () => {
    loadInitialServices();
    connectWebSocket();
});
