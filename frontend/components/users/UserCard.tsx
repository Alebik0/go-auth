import { UserData, usersApi } from "@/lib/api";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { useUserStore } from "@/app/contexts/UserContext";

import UserAvatar from "./UserAvatar";
import UserIdCard from "./UserIdCard";

const NAME_REGEX = /^[A-Za-z0-9\s]+$/;

function UserCard({ user }: { user: UserData }) {
  const userState = useUserStore((s) => s.state);
  const setUser = useUserStore((s) => s.setUser);

  const [name, setName] = useState(user.name);
  const [description, setDescription] = useState(user.description);
  const [nameError, setNameError] = useState<string | null>(null);
  const [descriptionError, setDescriptionError] = useState<string | null>(null);
  const [submitError, setSubmitError] = useState<string | null>(null);

  const router = useRouter();

  function validateName(name: string): boolean {
    if (!name.trim()) {
      setNameError("Name must not be empty.");
      return false;
    }

    if (name.length > 32) {
      setNameError("Name must not exceed 32 characters.");
      return false;
    }

    if (!NAME_REGEX.test(name)) {
      setNameError("Name may contain only letters, digits, and whitespace.");
      return false;
    }

    return true;
  }

  function validateDescription(description: string): boolean {
    if (description.length > 2048) {
      setDescriptionError("Description must not exceed 2048 characters.");
      return false;
    }

    return true;
  }

  function clearErrors() {
    setNameError(null);
    setDescriptionError(null);
    setSubmitError(null);
  }

  function onUserChange(newName: string, newDescription: string) {
    clearErrors();
    validateName(newName);
    validateDescription(newDescription);

    setName(newName);
    setDescription(newDescription);
  }

  function onNameChange(event: React.ChangeEvent<HTMLInputElement>) {
    const newName = event.target.value;
    const newDescription = description;

    onUserChange(newName, newDescription);
  }

  function onDescriptionChange(event: React.ChangeEvent<HTMLTextAreaElement>) {
    const newName = name;
    const newDescription = event.target.value;

    onUserChange(newName, newDescription);
  }

  function onBlur() {
    clearErrors();
    const user = userState.user;
    const validName = validateName(name);
    const validDescription = validateDescription(description);

    if (validName && validDescription && user != null) {
      usersApi
        .updateUser(user.id, {
          name: name,
          description: description,
        })
        .then((response) => setUser(response.data))
        .catch((error) => {
          switch (error.response?.status) {
            case 500:
              // Internal server error
              router.push("/error/500");
            default:
              setSubmitError("Failed update user");
          }
        });
    }
  }

  return (
    <div
      className="surface-container outline-variant border w-full max-w-2xl min-w-fit mx-auto p-8"
      style={{ borderRadius: "25px" }}
    >
      <div className="flex items-start gap-5">
        <UserAvatar>{user.name.charAt(0).toUpperCase()}</UserAvatar>

        <div className="flex-1">
          <UserIdCard id={user.id} />

          {/* Name */}
          <div className="mt-4">
            <label className="on-surface mb-1 block text-sm">Name</label>
            <input
              type="text"
              value={name}
              onChange={onNameChange}
              onBlur={() => onBlur()}
              maxLength={32}
              className={`${
                nameError == null
                  ? "surface-container-high on-surface-variant focus-surface-container-highest"
                  : "error-container on-error-container"
              } w-full rounded-xl px-4 py-2 text-white outline-none`}
            />
          </div>

          {/* Description */}
          <div className="mt-4">
            <label className="on-surface mb-1 block text-sm">Description</label>
            <textarea
              value={description}
              onChange={onDescriptionChange}
              onBlur={() => onBlur()}
              rows={4}
              maxLength={2048}
              className={`${
                descriptionError == null
                  ? "surface-container-high on-surface-variant focus-surface-container-highest"
                  : "error-container on-error-container"
              } w-full rounded-xl px-4 py-2 text-white outline-none`}
            />
          </div>

          {/* Submit error */}
          {submitError != null && (
            <p className="on-error-container mt-2 text-sm">{submitError}</p>
          )}
        </div>
      </div>
    </div>
  );
}

export default UserCard;
