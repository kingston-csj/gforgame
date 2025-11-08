using System;
using UnityEngine;
using UnityEngine.UI;

namespace Nova.Mvc
{
    public abstract class BaseView : MonoBehaviour
    {
        public bool visible
        {
            get => gameObject.activeSelf;
            set => gameObject.SetActive(value);
        }

        protected void RegisterClickEvent(GameObject node, Action callback)
        {
            if (node.TryGetComponent(out Button button))
            {
                button.onClick.AddListener(() => callback?.Invoke());
            }
            else if (node.TryGetComponent(out Toggle toggle))
            {
                toggle.onValueChanged.AddListener((value) => callback?.Invoke());
            }
        }
    }
    
}