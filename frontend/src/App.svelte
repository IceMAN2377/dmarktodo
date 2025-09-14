<script>
  import logo from './assets/images/logo-universal.png'
  import { GetTasks, AddTask, DeleteTask } from '../wailsjs/go/backend/App.js';

  // Define the Task interface
  let tasks = [];
  let newTaskTitle = "";
  let showModal = false;
  let modalMessage = "";
  let modalType = "error";

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

  // Function to get all tasks
  async function getTasks() {
    try {
      const fetchedTasks = await GetTasks();
      console.log("Loaded tasks:", fetchedTasks);

      // Принудительно обновляем массив задач
      tasks = [...fetchedTasks];

      // Если массив пуст, убедимся, что это действительно пустой массив
      if (tasks.length === 0) {
        tasks = [];
      }
    } catch (error) {
      console.error("Error fetching tasks:", error);
      tasks = [];

    }
  }

  // Function to add a new task
  async function addTask() {
    // Проверяем на пустую строку
    if (newTaskTitle.trim() === "") {
      showAlert("Пожалуйста, введите название задачи!", "warning");
      return;
    }

    try {
      const newTask = await AddTask(newTaskTitle);
      console.log("Added task:", newTask);

      // Перезагружаем все задачи после добавления
      await getTasks();


      newTaskTitle = "";
    } catch (error) {
      console.error("Error adding task:", error);
      showAlert("Ошибка при добавлении задачи", "error");
    }
  }

  // Function to delete a task
  async function deleteTask(id) {
    try {
      const success = await DeleteTask(id);
      console.log("Delete result:", success, "for task ID:", id);

      if (success) {
        tasks = tasks.filter(task => task.id !== id);
        if (tasks.length === 0) {
          tasks = [];
        }
        // Перезагружаем все задачи после удаления
        await getTasks();
        showAlert("Задача успешно удалена!", "success");
      } else {
        console.error("Failed to delete task with ID:", id);
        showAlert("Не удалось удалить задачу", "error");
      }
    } catch (error) {
      console.error("Error deleting task:", error);
      showAlert("Ошибка при удалении задачи", "error");
    }
  }

  // Handle keydown for modal (close on Escape)
  function handleKeydown(event) {
    if (event.key === 'Escape' && showModal) {
      closeModal();
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
              on:keypress={(e) => e.key === 'Enter' && addTask()}
              class="input"
      />
      <button class="btn" on:click={addTask}>Add Task</button>
    </div>

    <div class="task-list">
      {#if tasks.length === 0}
        <p class="no-tasks">No tasks yet. Add one above!</p>
      {:else}
        <ul>
          {#each tasks as task (task.id)}
            <li class="task-item">
              <span class="task-id">#{task.id}</span>
              <span class="task-title">{task.title}</span>
              <span class="task-status status-{task.status}">{task.status}</span>
              <span class="task-priority priority-{task.priority}">{task.priority}</span>
              <button class="delete-btn" on:click={() => deleteTask(task.id)}>Delete</button>
            </li>
          {/each}
        </ul>
      {/if}
    </div>

    <!-- Debug info -->
    <div class="debug-info">
      <p>Total tasks: {tasks.length}</p>
      <button class="btn" on:click={getTasks}>Refresh Tasks</button>
    </div>
  </div>

  <!-- Modal -->
  {#if showModal}
    <div class="modal-overlay" on:click={closeModal}>
      <div class="modal modal-{modalType}" on:click|stopPropagation>
        <div class="modal-header">
          <h3>{modalType === 'error' ? '⚠️ Ошибка' : modalType === 'warning' ? '⚠️ Внимание' : '✅ Успешно'}</h3>
          <button class="modal-close" on:click={closeModal}>×</button>
        </div>
        <div class="modal-body">
          <p>{modalMessage}</p>
        </div>
        <div class="modal-footer">
          <button class="btn-modal" on:click={closeModal}>ОК</button>
        </div>
      </div>
    </div>
  {/if}
</main>

<style>
  main {
    max-width: 800px;
    margin: 0 auto;
    padding: 20px;
  }

  .todo-app {
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 100%;
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

  .modal-footer {
    padding: 10px 20px 20px 20px;
    text-align: right;
  }

  .btn-modal {
    background-color: #4a90e2;
    color: white;
    border: none;
    border-radius: 4px;
    padding: 8px 20px;
    cursor: pointer;
    font-size: 14px;
    font-weight: 500;
  }

  .btn-modal:hover {
    background-color: #357abd;
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
    color: #333;
    margin-bottom: 20px;
    text-align: center;
  }

  .add-task {
    display: flex;
    width: 100%;
    margin-bottom: 20px;
  }

  .input {
    flex: 1;
    border: 1px solid #ddd;
    border-radius: 4px;
    outline: none;
    height: 40px;
    line-height: 40px;
    padding: 0 15px;
    font-size: 16px;
    background-color: #f9f9f9;
  }

  .input:focus {
    border-color: #4a90e2;
    background-color: #fff;
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
    background-color: #357abD;
  }

  .task-list {
    width: 100%;
  }

  .no-tasks {
    text-align: center;
    color: #888;
    font-style: italic;
  }

  ul {
    list-style-type: none;
    padding: 0;
    width: 100%;
  }

  .task-item {
    display: flex;
    align-items: center;
    padding: 12px 15px;
    background-color: #f9f9f9;
    border-radius: 4px;
    margin-bottom: 10px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  }

  .task-id {
    font-size: 14px;
    color: #888;
    margin-right: 10px;
    min-width: 30px;
  }

  .task-title {
    flex: 1;
    font-size: 16px;
    color: #333;
    margin-right: 10px;
  }

  .task-status {
    font-size: 12px;
    padding: 2px 8px;
    border-radius: 12px;
    margin-right: 10px;
    text-transform: uppercase;
    font-weight: bold;
  }

  .status-active {
    background-color: #e3f2fd;
    color: #1976d2;
  }

  .status-done {
    background-color: #e8f5e8;
    color: #388e3c;
  }

  .task-priority {
    font-size: 12px;
    padding: 2px 8px;
    border-radius: 12px;
    margin-right: 10px;
    text-transform: uppercase;
    font-weight: bold;
  }

  .priority-low {
    background-color: #f3e5f5;
    color: #7b1fa2;
  }

  .priority-medium {
    background-color: #fff3e0;
    color: #f57c00;
  }

  .priority-high {
    background-color: #ffebee;
    color: #d32f2f;
  }

  .delete-btn {
    background-color: #ff4d4d;
    color: white;
    border: none;
    border-radius: 4px;
    padding: 5px 10px;
    cursor: pointer;
    font-size: 14px;
  }

  .delete-btn:hover {
    background-color: #e60000;
  }

  .debug-info {
    margin-top: 20px;
    padding: 10px;
    background-color: #f0f0f0;
    border-radius: 4px;
    color: #333;
  }

  .debug-info p {
    margin: 5px 0;
    font-size: 14px;
  }
</style>