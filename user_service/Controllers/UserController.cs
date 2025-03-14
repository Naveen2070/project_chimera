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

            var user = await _userService.CreateAsync(userCreateDto);
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
