<script>
  import logo from './assets/images/logo-universal.png'

  
  // Define the Task interface
  let tasks = [];
  let newTaskTitle = "";
  
  // Function to get all tasks
  async function getTasks() {
    try {
      // Call the Go function directly using the window object
      tasks = await window.go.backend.App.GetTasks();
    } catch (error) {
      console.error("Error fetching tasks:", error);
    }
  }
  
  // Function to add a new task
  async function addTask() {
    if (newTaskTitle.trim() === "") return;
    
    try {
      // Call the Go function directly using the window object
      const newTask = await window.go.backend.App.AddTask(newTaskTitle);
      tasks = [...tasks, newTask];
      newTaskTitle = "";
    } catch (error) {
      console.error("Error adding task:", error);
    }
  }
  
  // Function to delete a task
  async function deleteTask(id) {
    try {
      // Call the Go function directly using the window object
      const success = await window.go.backend.App.DeleteTask(id);
      if (success) {
        // Remove the task from the local array
        tasks = tasks.filter(task => task.id !== id);
      } else {
        console.error("Failed to delete task with ID:", id);
      }
    } catch (error) {
      console.error("Error deleting task:", error);
    }
  }
  
  // Load tasks when the component is mounted
  getTasks();
</script>

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
              <button class="delete-btn" on:click={() => deleteTask(task.id)}>Delete</button>
            </li>
          {/each}
        </ul>
      {/if}
    </div>
  </div>
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
</style>
