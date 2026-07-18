using UnityEngine;

namespace Nova.Ui
{
    public interface IUiFactory
    {
        Transform GetNodeByLayer(LayerIds layer);
    }
}