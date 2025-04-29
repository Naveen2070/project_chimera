//	Copyright 2025 Naveen R
//
//		Licensed under the Apache License, Version 2.0 (the "License");
//		you may not use this file except in compliance with the License.
//		You may obtain a copy of the License at
//
//		http://www.apache.org/licenses/LICENSE-2.0
//
//		Unless required by applicable law or agreed to in writing, software
//		distributed under the License is distributed on an "AS IS" BASIS,
//		WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//		See the License for the specific language governing permissions and
//		limitations under the License.
using Microsoft.AspNetCore.Mvc;
using user_service.Database;
using user_service.Model;
using user_service.Services.Interfaces;

namespace user_service.Controllers
{
    [Route("/")]
    [ApiController]
    public class UserController : ControllerBase
    {
        private readonly IUserService _userService;

        public UserController(IUserService userService)
        {
            _userService = userService;
        }

        // GET /
        [HttpGet]
        public async Task<IActionResult> GetAllUsers()
        {
            var users = await _userService.GetAllAsync();
            return Ok(users);
        }

        // GET /{id}
        [HttpGet("{id}")]
        public async Task<IActionResult> GetUserById(Guid id)
        {
            var user = await _userService.GetByIdAsync(id);
            if (user == null)
                return NotFound($"User with id {id} not found.");

            return Ok(user);
        }

        // POST /
        [HttpPost]
        public async Task<IActionResult> CreateUser([FromBody] UserCreateDTO userCreateDto)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var user = new UserDTO();

            try
            {
                user = await _userService.CreateAsync(userCreateDto);
            }
            catch (Exception ex)
            { 
                return BadRequest(ex.Message);
            }
            return CreatedAtAction(nameof(GetUserById), new { id = user.Id }, user);
        }

        // PUT /{id}
        [HttpPut("{id}")]
        public async Task<IActionResult> UpdateUser(Guid id, [FromBody] UserUpdateDTO userUpdateDto)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var updatedUser = await _userService.UpdateAsync(id, userUpdateDto);
            if (updatedUser == null)
                return NotFound($"User with id {id} not found.");

            return Ok(updatedUser);
        }

        // PUT /credentials/{id}
        [HttpPut("credentials/{id}")]
        public async Task<IActionResult> UpdateUserCredentials(Guid id, [FromBody] UserCredentialsUpdateDTO credentialsUpdateDto)
        {
            if (!ModelState.IsValid)
                return BadRequest(ModelState);

            var result = await _userService.UpdateCredentialsAsync(id, credentialsUpdateDto);
            if (!result)
                return NotFound($"User with id {id} not found or update failed.");

            return Ok($"User credentials updated successfully.");
        }

        // PATCH /{id}
        [HttpPatch("{id}")]
        public async Task<IActionResult> DeleteUser(Guid id)
        {
            var softDeleted = await _userService.SoftDeleteAsync(id);
            if (!softDeleted)
                return NotFound($"User with id {id} not found.");

            return NoContent();
        }

        // PATCH /erase/{id}
        [HttpDelete("erase/{id}")]
        public async Task<IActionResult> SoftDeleteUser(Guid id)
        {

            var deleted = await _userService.DeleteAsync(id);
            if (!deleted)
                return NotFound($"User with id {id} not found.");

            return Ok($"User with id {id} successfully deleted from database.");
        }
    }
}
