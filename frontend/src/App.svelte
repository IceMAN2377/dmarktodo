<script>
  import logo from './assets/images/logo-universal.png'
  import { GetTasks, AddTask, DeleteTask, ToggleStatus } from '../wailsjs/go/backend/App.js';

  // Define the Task interface
  let tasks = [];
  let newTaskTitle = "";
  let showModal = false;
  let modalMessage = "";
  let modalType = "error";

  // Priority selection state
  let showPriorityModal = false;
  let selectedPriority = "medium";

  // Confirmation modal state
  let showConfirmModal = false;
  let taskToDelete = null;

  // Filter state
  let currentFilter = "all"; // "all", "active", "done"

  // Sort state
  let sortByNewest = true; // true = newest first (DESC), false = oldest first (ASC)

  // Computed properties for active and completed tasks
  $: activeTasks = tasks.filter(task => task.status === 'active');
  $: completedTasks = tasks.filter(task => task.status === 'done');

  // Computed property for filtered tasks based on current filter
  $: filteredTasks = currentFilter === "all" ? tasks :
                     currentFilter === "active" ? activeTasks : completedTasks;

  // Function to show modal
  function showAlert(message, type = "error") {
    modalMessage = message;
    modalType = type;
    showModal = true;
  }

  // Function to close modal
  function closeModal() {
    showModal = false;
    modalMessage = "";
  }

  // Function to show priority modal
  function showPrioritySelection() {
    showPriorityModal = true;
  }

  // Function to close priority modal
  function closePriorityModal() {
    showPriorityModal = false;
    selectedPriority = "medium";
  }

  // Function to show confirmation modal
  function showDeleteConfirmation(task) {
    taskToDelete = task;
    showConfirmModal = true;
  }

  // Function to close confirmation modal
  function closeConfirmModal() {
    showConfirmModal = false;
    taskToDelete = null;
  }

  // Function to get all tasks
  async function getTasks() {
    try {
      const fetchedTasks = await GetTasks(sortByNewest);
      console.log("Loaded tasks:", fetchedTasks);

      // Force update the tasks array
      tasks = [...fetchedTasks];

      // If array is empty, make sure it's actually an empty array
      if (tasks.length === 0) {
        tasks = [];
      }
    } catch (error) {
      console.error("Error fetching tasks:", error);
      tasks = [];
    }
  }

  // Function to toggle sort order
  async function toggleSort() {
    sortByNewest = !sortByNewest;
    await getTasks();
  }

  // Function to start adding a task
  async function startAddTask() {
    // Check for empty string
    if (newTaskTitle.trim() === "") {
      showAlert("Please enter a task title!", "warning");
      return;
    }

    // Show priority selection modal
    showPrioritySelection();
  }

  // Function to add a new task with priority
  async function addTaskWithPriority() {
    try {
      const newTask = await AddTask(newTaskTitle, selectedPriority);
      console.log("Added task:", newTask);

      // Reload all tasks after adding
      await getTasks();

      newTaskTitle = "";
      closePriorityModal();
      showAlert("Task added successfully!", "success");
    } catch (error) {
      console.error("Error adding task:", error);
      showAlert("Error adding task", "error");
      closePriorityModal();
    }
  }

  // Function to confirm and delete a task
  async function confirmDelete() {
    if (!taskToDelete) return;

    try {
      const success = await DeleteTask(taskToDelete.id);
      console.log("Delete result:", success, "for task ID:", taskToDelete.id);

      if (success) {
        await getTasks();
        showAlert("Task deleted successfully!", "success");
      } else {
        console.error("Failed to delete task with ID:", taskToDelete.id);
        showAlert("Failed to delete task", "error");
      }
    } catch (error) {
      console.error("Error deleting task:", error);
      showAlert("Error deleting task", "error");
    }

    closeConfirmModal();
  }

  // Function to toggle task status
  async function toggleTaskStatus(id, Status) {
    try {
      console.log("Toggling status for task:", id, "current status:", Status);

      const updatedTask = await ToggleStatus(id, Status);
      console.log("Task status updated:", updatedTask);

      // Reload all tasks after status change
      await getTasks();

      const statusText = Status === 'active' ? 'completed' : 'active';
      showAlert(`Task marked as ${statusText}!`, "success");
    } catch (error) {
      console.error("Error toggling task status:", error);
      showAlert("Error changing task status", "error");
    }
  }

  // Handle keydown for modal (close on Escape)
  function handleKeydown(event) {
    if (event.key === 'Escape') {
      if (showModal) {
        closeModal();
      }
      if (showPriorityModal) {
        closePriorityModal();
      }
      if (showConfirmModal) {
        closeConfirmModal();
      }
    }
  }

  // Load tasks when the component is mounted
  getTasks();
</script>

<svelte:window on:keydown={handleKeydown} />

<main>
  <div class="todo-app">
    <img alt="Todo App logo" id="logo" src="{logo}">
    <h1>Todo List</h1>

    <div class="add-task">
      <input
              type="text"
              placeholder="Enter a new task..."
              bind:value={newTaskTitle}
              on:keypress={(e) => e.key === 'Enter' && startAddTask()}
              class="input"
      />
      <button class="btn" on:click={startAddTask}>Add Task</button>
    </div>

    <!-- Filter buttons -->
    <div class="filter-buttons">
      <button
        class="filter-btn {currentFilter === 'all' ? 'active' : ''}"
        on:click={() => currentFilter = 'all'}
      >
        All ({tasks.length})
      </button>
      <button
        class="filter-btn {currentFilter === 'active' ? 'active' : ''}"
        on:click={() => currentFilter = 'active'}
      >
        Active ({activeTasks.length})
      </button>
      <button
        class="filter-btn {currentFilter === 'done' ? 'active' : ''}"
        on:click={() => currentFilter = 'done'}
      >
        Completed ({completedTasks.length})
      </button>
    </div>

    <!-- Sort buttons -->
    <div class="sort-buttons">
      <button
        class="sort-btn {sortByNewest ? 'active' : ''}"
        on:click={toggleSort}
        title="Sort by creation time"
      >
        {sortByNewest ? '🔽 Newest First' : '🔼 Oldest First'}
      </button>
    </div>

    <!-- Tasks list -->
    <div class="tasks-container">
      <div class="tasks-section">
        <div class="task-list">
          {#if filteredTasks.length === 0}
            <p class="no-tasks">
              {#if currentFilter === 'all'}
                No tasks
              {:else if currentFilter === 'active'}
                No active tasks
              {:else}
                No completed tasks
              {/if}
            </p>
          {:else}
            <ul>
              {#each filteredTasks as task (task.id)}
                <li class="task-item {task.status === 'active' ? 'task-active' : 'task-completed'}">
                  <div class="task-info">
                    <span class="task-id">#{task.id}</span>
                    <span class="task-title {task.status === 'done' ? 'completed-title' : ''}">{task.title}</span>
                    <span class="task-priority priority-{task.priority}">{task.priority}</span>
                    {#if task.completed_at && task.status === 'done'}
                      <span class="completed-date">
                        Completed: {new Date(task.completed_at).toLocaleString('en-US')}
                      </span>
                    {/if}
                  </div>
                  <div class="task-actions">
                    {#if task.status === 'active'}
                      <button
                        class="complete-btn"
                        on:click={() => toggleTaskStatus(task.id, task.status)}
                        title="Mark as completed"
                      >
                        ✓
                      </button>
                    {:else}
                      <button
                        class="reactivate-btn"
                        on:click={() => toggleTaskStatus(task.id, task.status)}
                        title="Mark as active"
                      >
                        ↶
                      </button>
                    {/if}
                    <button
                      class="delete-btn"
                      on:click={() => showDeleteConfirmation(task)}
                      title="Delete task"
                    >
                      ✕
                    </button>
                  </div>
                </li>
              {/each}
            </ul>
          {/if}
        </div>
      </div>
    </div>

  <!-- Priority Selection Modal -->
  {#if showPriorityModal}
    <div class="modal-overlay" on:click={closePriorityModal}>
      <div class="modal priority-modal" on:click|stopPropagation>
        <div class="modal-header">
          <h3>Select Task Priority</h3>
          <button class="modal-close" on:click={closePriorityModal}>×</button>
        </div>
        <div class="modal-body">
          <div class="priority-options">
            <label class="priority-option">
              <input type="radio" bind:group={selectedPriority} value="low" />
              <span class="priority-label priority-low">Low</span>
            </label>
            <label class="priority-option">
              <input type="radio" bind:group={selectedPriority} value="medium" />
              <span class="priority-label priority-medium">Medium</span>
            </label>
            <label class="priority-option">
              <input type="radio" bind:group={selectedPriority} value="high" />
              <span class="priority-label priority-high">High</span>
            </label>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn-modal btn-cancel" on:click={closePriorityModal}>Cancel</button>
          <button class="btn-modal btn-confirm" on:click={addTaskWithPriority}>Add Task</button>
        </div>
      </div>
    </div>
  {/if}

  <!-- Delete Confirmation Modal -->
  {#if showConfirmModal && taskToDelete}
    <div class="modal-overlay" on:click={closeConfirmModal}>
      <div class="modal modal-warning" on:click|stopPropagation>
        <div class="modal-header">
          <h3>🗑️ Confirm Deletion</h3>
          <button class="modal-close" on:click={closeConfirmModal}>×</button>
        </div>
        <div class="modal-body">
          <p>Are you sure you want to delete the task "{taskToDelete.title}"?</p>
          <p style="color: #999; font-size: 14px; margin-top: 10px;">This action cannot be undone.</p>
        </div>
        <div class="modal-footer">
          <button class="btn-modal btn-cancel" on:click={closeConfirmModal}>Cancel</button>
          <button class="btn-modal btn-delete" on:click={confirmDelete}>Delete</button>
        </div>
      </div>
    </div>
  {/if}

  <!-- Alert Modal -->
  {#if showModal}
    <div class="modal-overlay" on:click={closeModal}>
      <div class="modal modal-{modalType}" on:click|stopPropagation>
        <div class="modal-header">
          <h3>{modalType === 'error' ? '⚠️ Error' : modalType === 'warning' ? '⚠️ Warning' : '✅ Success'}</h3>
          <button class="modal-close" on:click={closeModal}>×</button>
        </div>
        <div class="modal-body">
          <p>{modalMessage}</p>
        </div>
        <div class="modal-footer">
          <button class="btn-modal" on:click={closeModal}>OK</button>
        </div>
      </div>
    </div>
  {/if}
</main>

<style>
  main {
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
  }

  .todo-app {
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 100%;
  }

  .filter-buttons {
    display: flex;
    gap: 10px;
    margin: 20px 0;
  }

  .filter-btn {
    padding: 8px 16px;
    border: 2px solid rgba(255, 255, 255, 0.2);
    background-color: rgba(255, 255, 255, 0.1);
    color: rgba(255, 255, 255, 0.8);
    border-radius: 20px;
    cursor: pointer;
    font-size: 14px;
    transition: all 0.2s ease;
  }

  .filter-btn:hover {
    background-color: rgba(255, 255, 255, 0.15);
    border-color: rgba(255, 255, 255, 0.3);
  }

  .filter-btn.active {
    background-color: #4a90e2;
    border-color: #4a90e2;
    color: white;
  }

  .sort-buttons {
    display: flex;
    gap: 10px;
    margin: 10px 0;
  }

  .sort-btn {
    padding: 8px 16px;
    border: 2px solid rgba(255, 255, 255, 0.2);
    background-color: rgba(255, 255, 255, 0.1);
    color: rgba(255, 255, 255, 0.8);
    border-radius: 20px;
    cursor: pointer;
    font-size: 14px;
    transition: all 0.2s ease;
  }

  .sort-btn:hover {
    background-color: rgba(255, 255, 255, 0.15);
    border-color: rgba(255, 255, 255, 0.3);
  }

  .sort-btn.active {
    background-color: #28a745;
    border-color: #28a745;
    color: white;
  }

  .tasks-container {
    width: 100%;
    margin-top: 10px;
  }

  .tasks-section {
    background-color: rgba(255, 255, 255, 0.05);
    border-radius: 8px;
    padding: 20px;
    min-height: 300px;
  }

  .tasks-section h2 {
    margin-top: 0;
    margin-bottom: 15px;
    color: #fff;
    font-size: 18px;
    text-align: center;
    border-bottom: 2px solid rgba(255, 255, 255, 0.1);
    padding-bottom: 10px;
  }

  /* Modal styles */
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    animation: fadeIn 0.2s ease-out;
  }

  .modal {
    background: white;
    border-radius: 8px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
    width: 90%;
    max-width: 400px;
    animation: slideUp 0.3s ease-out;
  }

  .priority-modal {
    max-width: 450px;
  }

  .modal-error {
    border-top: 4px solid #f44336;
  }

  .modal-warning {
    border-top: 4px solid #ff9800;
  }

  .modal-success {
    border-top: 4px solid #4caf50;
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px 20px 10px 20px;
    border-bottom: 1px solid #eee;
  }

  .modal-header h3 {
    margin: 0;
    color: #333;
    font-size: 18px;
  }

  .modal-close {
    background: none;
    border: none;
    font-size: 24px;
    cursor: pointer;
    color: #999;
    padding: 0;
    width: 30px;
    height: 30px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
  }

  .modal-close:hover {
    background-color: #f5f5f5;
    color: #666;
  }

  .modal-body {
    padding: 20px;
  }

  .modal-body p {
    margin: 0;
    color: #666;
    font-size: 16px;
    line-height: 1.5;
  }

  .priority-options {
    display: flex;
    flex-direction: column;
    gap: 15px;
  }

  .priority-option {
    display: flex;
    align-items: center;
    cursor: pointer;
    padding: 10px;
    border: 1px solid #ddd;
    border-radius: 6px;
    transition: all 0.2s ease;
  }

  .priority-option:hover {
    background-color: #f8f9fa;
    border-color: #ccc;
  }

  .priority-option input[type="radio"] {
    margin-right: 10px;
    cursor: pointer;
  }

  .priority-label {
    font-weight: 500;
    padding: 4px 10px;
    border-radius: 12px;
    text-transform: uppercase;
    font-size: 12px;
    color: white;
  }

  .priority-low {
    background-color: #9c27b0;
  }

  .priority-medium {
    background-color: #ff9800;
  }

  .priority-high {
    background-color: #f44336;
  }

  .modal-footer {
    padding: 10px 20px 20px 20px;
    text-align: right;
    display: flex;
    gap: 10px;
    justify-content: flex-end;
  }

  .btn-modal {
    border: none;
    border-radius: 4px;
    padding: 8px 20px;
    cursor: pointer;
    font-size: 14px;
    font-weight: 500;
  }

  .btn-cancel {
    background-color: #6c757d;
    color: white;
  }

  .btn-cancel:hover {
    background-color: #5a6268;
  }

  .btn-confirm {
    background-color: #4a90e2;
    color: white;
  }

  .btn-confirm:hover {
    background-color: #357abd;
  }

  .btn-delete {
    background-color: #dc3545;
    color: white;
  }

  .btn-delete:hover {
    background-color: #c82333;
  }

  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }

  @keyframes slideUp {
    from {
      opacity: 0;
      transform: translateY(30px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  #logo {
    display: block;
    width: 120px;
    height: auto;
    margin: 0 auto 20px;
  }

  h1 {
    color: #fff;
    margin-bottom: 20px;
    text-align: center;
  }

  .add-task {
    display: flex;
    width: 100%;
    max-width: 500px;
    margin-bottom: 20px;
  }

  .input {
    flex: 1;
    border: 1px solid rgba(255, 255, 255, 0.2);
    border-radius: 4px;
    outline: none;
    height: 40px;
    line-height: 40px;
    padding: 0 15px;
    font-size: 16px;
    background-color: rgba(255, 255, 255, 0.1);
    color: #fff;
  }

  .input::placeholder {
    color: rgba(255, 255, 255, 0.6);
  }

  .input:focus {
    border-color: #4a90e2;
    background-color: rgba(255, 255, 255, 0.15);
  }

  .btn {
    height: 40px;
    padding: 0 15px;
    margin-left: 10px;
    background-color: #4a90e2;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 16px;
  }

  .btn:hover {
    background-color: #357abd;
  }

  .task-list {
    width: 100%;
  }

  .no-tasks {
    text-align: center;
    color: rgba(255, 255, 255, 0.6);
    font-style: italic;
    margin-top: 20px;
  }

  ul {
    list-style-type: none;
    padding: 0;
    width: 100%;
  }

  .task-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 15px;
    border-radius: 8px;
    margin-bottom: 10px;
    border: 1px solid transparent;
    transition: all 0.2s ease;
  }

  .task-active {
    background-color: rgba(74, 144, 226, 0.1);
    border-color: rgba(74, 144, 226, 0.3);
  }

  .task-completed {
    background-color: rgba(76, 175, 80, 0.1);
    border-color: rgba(76, 175, 80, 0.3);
  }

  .task-info {
    display: flex;
    align-items: center;
    flex: 1;
    gap: 10px;
  }

  .task-actions {
    display: flex;
    gap: 8px;
  }

  .task-id {
    font-size: 12px;
    color: rgba(255, 255, 255, 0.6);
    min-width: 30px;
  }

  .task-title {
    font-size: 16px;
    color: #fff;
    flex: 1;
  }

  .completed-title {
    text-decoration: line-through;
    color: rgba(255, 255, 255, 0.7);
  }

  .task-priority {
    font-size: 12px;
    padding: 3px 8px;
    border-radius: 12px;
    text-transform: uppercase;
    font-weight: bold;
  }

  .priority-low {
    background-color: rgba(156, 39, 176, 0.2);
    color: #ce93d8;
  }

  .priority-medium {
    background-color: rgba(255, 152, 0, 0.2);
    color: #ffcc02;
  }

  .priority-high {
    background-color: rgba(244, 67, 54, 0.2);
    color: #ff8a80;
  }

  .completed-date {
    font-size: 12px;
    color: rgba(255, 255, 255, 0.5);
  }

  .complete-btn {
    background-color: #4caf50;
    color: white;
    border: none;
    border-radius: 50%;
    width: 32px;
    height: 32px;
    cursor: pointer;
    font-size: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background-color 0.2s;
  }

  .complete-btn:hover {
    background-color: #45a049;
  }

  .reactivate-btn {
    background-color: #ff9800;
    color: white;
    border: none;
    border-radius: 50%;
    width: 32px;
    height: 32px;
    cursor: pointer;
    font-size: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background-color 0.2s;
  }

  .reactivate-btn:hover {
    background-color: #f57c00;
  }

  .delete-btn {
    background-color: #f44336;
    color: white;
    border: none;
    border-radius: 50%;
    width: 32px;
    height: 32px;
    cursor: pointer;
    font-size: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background-color 0.2s;
  }

  .delete-btn:hover {
    background-color: #e53935;
  }

  .debug-info p {
    margin: 5px 0;
    font-size: 14px;
  }

  /* Responsive design */
  @media (max-width: 768px) {
    .tasks-container {
      grid-template-columns: 1fr;
      gap: 20px;
    }

    .task-item {
      flex-direction: column;
      align-items: flex-start;
      gap: 10px;
    }

    .task-actions {
      align-self: flex-end;
    }
  }
  </style>